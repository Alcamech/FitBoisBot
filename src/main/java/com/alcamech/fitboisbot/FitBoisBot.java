package com.alcamech.fitboisbot;

import com.alcamech.fitboisbot.model.FitBoiGg;
import com.alcamech.fitboisbot.model.FitBoiRecord;
import com.alcamech.fitboisbot.model.FitBoiUser;
import com.alcamech.fitboisbot.respository.RecordRepository;
import com.alcamech.fitboisbot.respository.GgRepository;
import com.alcamech.fitboisbot.respository.UserRepository;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;
import org.telegram.telegrambots.bots.TelegramLongPollingBot;
import org.telegram.telegrambots.meta.api.methods.send.SendMessage;
import org.telegram.telegrambots.meta.api.objects.Message;
import org.telegram.telegrambots.meta.api.objects.Update;
import org.telegram.telegrambots.meta.api.objects.User;
import org.telegram.telegrambots.meta.exceptions.TelegramApiException;

import java.time.LocalDate;
import java.time.Year;
import java.time.ZoneId;
import java.time.ZonedDateTime;
import java.time.format.DateTimeFormatter;
import java.util.*;
import java.util.stream.Collectors;

import static com.alcamech.fitboisbot.Constants.*;

@Component
public class FitBoisBot extends TelegramLongPollingBot {
    private static final Logger logger = LoggerFactory.getLogger(FitBoisBot.class);

    @Autowired
    RecordRepository recordRepository;

    @Autowired
    UserRepository userRepository;

    @Autowired
    GgRepository ggRepository;

    @Value("${token}")
    private String token;

    @Value("${username}")
    private String username;

    @Override
    public String getBotUsername() {
        return username;
    }

    @Override
    public String getBotToken() {
        return token;
    }

    private Long lastActivityPostUserId;
    private boolean isFastestGGAvailable;

    @Scheduled(cron = "0 0 9 1 * ?", zone = "America/New_York")
    public void awardMostActiveUserForPreviousMonth() {
        ZonedDateTime nowInEST = ZonedDateTime.now(ZoneId.of("America/New_York"));
        System.out.println("SCHEDULED - Time now in EST: " + nowInEST);

        String year = String.valueOf(nowInEST.getYear());
        String month = String.format("%02d", nowInEST.getMonthValue() - 1);
        if(month.equals("00")) {
            year = String.valueOf(Integer.parseInt(year) - 1); // Decrement year if month is December of the previous year
            month = "12";
        }

        Long maxActivityCount = recordRepository.findMaxActivityCountByYearAndMonth(month);
        if (maxActivityCount != null && maxActivityCount > 0) {
            List<Long> mostActiveUserIds = recordRepository.findAllUsersWithMaxCount(month, maxActivityCount);

            String resetMessage = "Monthly counts have been reset";
            FitBoiUser anActiveUser = userRepository.findById(mostActiveUserIds.get(0)).orElse(null);
            sendText(anActiveUser.getGroupId(), resetMessage);

            for (Long userId : mostActiveUserIds) {
                FitBoiUser mostActiveUser = userRepository.findById(userId).orElse(null);
                if (mostActiveUser != null) {
                    String congratsMessage = "Congratulations \u2B50 " + mostActiveUser.getName() + " \u2B50 for being the most active user for " + month + "/" + year + " \uD83C\uDFC6";
                    sendText(mostActiveUser.getGroupId(), congratsMessage);
                    String rewardMessage = "Here's your reward. \uD83D\uDCB0 You've won 100 FitBoi Tokens! \uD83D\uDCB0";
                    sendText(mostActiveUser.getGroupId(), rewardMessage);
                }
            }
        }
    }

    @Override
    public void onUpdateReceived(Update update) {
        try {
            Message msg = getMessageFromUpdate(update);
            FitBoiUser user = getUserFromMessage(msg);
            System.out.println("Message: " + msg);
            System.out.println("User: " + user);
            Long chatId = msg.getChat().getId();

            isFastestGG(update, msg, chatId, user);

            if (msg.hasPhoto())
                handlePhotoMessage(msg, chatId);

            if (msg.isCommand())
                handleCommandMessage(msg, chatId);

        } catch (Exception e) {
            logger.error("An error occurred in onUpdateReceived method: ", e);
        }
    }

    private Message getMessageFromUpdate(Update update) {
        if (update.getEditedMessage() != null) {
            return update.getEditedMessage();
        } else {
            Message msg = update.getMessage();
            if (msg.hasPhoto()) {
                lastActivityPostUserId = msg.getFrom().getId();
                isFastestGGAvailable = true;
            }
            return msg;
        }
    }

    private void handlePhotoMessage(Message msg, Long chatId) {
        if (msg.getCaption() != null) {
            String msgCaption = msg.getCaption();
            Long userId = msg.getFrom().getId();
            Map<String, String> captionContents = parseMessageCaption(msgCaption, chatId);
            getRecordFromCaption(userId, captionContents);

            String totalActivitiesMessage = getActivityCountsMessage();
            sendText(chatId, totalActivitiesMessage);
        }
    }

    private void handleCommandMessage(Message msg, Long chatId) {
        String commandText = msg.getText();
        if (commandText.equals(Commands.HELP.toString())) {
            onHelp(chatId);
        }
        if (commandText.equals(Commands.FASTGG.toString())) {
            onFastGG(chatId);
        }
    }

    /**
    * Response for the /help command
    *
    * @param chatId chat id to send response to
    */
    public void onHelp(Long chatId) {
        String message = "Nothing to see here... yet";
        sendText(chatId, message);
    }

    /**
    * Response for the /fastgg commandALTER TABLE fit_boi_user DROP COLUMN fast_gg_count;

    *
    * @param chatId chat id to send response to
    */
    public void onFastGG(Long chatId) {
        String message = getFastestGgCountsMessage(chatId);
        sendText(chatId, message);
    }

    /**
    * Sends a message
    *
    * @param chatId chat id to send message to
    * @param messageContent content of the message
    */
    public void sendText(Long chatId, String messageContent){
        SendMessage sm = SendMessage.builder()
                .chatId(chatId.toString())
                .text(messageContent).build();
        try {
            execute(sm);
        } catch (TelegramApiException e) {
            throw new RuntimeException(e);
        }
    }

    /**
    * Sends a message as a reply
    *
    * @param chatId chat id to send message to
    * @param messageContent content of the message
    * @param messageId message id to reply to
    */
    public void sendTextAsReply(Long chatId, String messageContent, int messageId){
        SendMessage sm = SendMessage.builder()
                .chatId(chatId.toString())
                .text(messageContent)
                .replyToMessageId(messageId)
                .build();
        try {
            execute(sm);
        } catch (TelegramApiException e) {
            throw new RuntimeException(e);
        }
    }

    /**
    * Determines if the message is the fastest gg. We only want to reply to the
    * first gg after the last activity post. If the message is the fastest gg
    * update the count in the database.
    *
    * @param update update
    * @param msg message
    * @param chatId chat id
    */
    public void isFastestGG(Update update, Message msg, Long chatId, FitBoiUser user) {
        String fastestGG = "Fastest gg in the west";
        if (!isFastestGGAvailable) {
            return;
        }

        // update is an edit from activity poster
        if (update.getEditedMessage() != null && update.getEditedMessage().getFrom().getId().equals(lastActivityPostUserId)) {
            return;
        }

        if (isGG(msg.getText())) {
            int currentYear = LocalDate.now().getYear();
            System.out.println("Current Year: " + currentYear);
            ggRepository.updateGgCountForCurrentYear(user.getId(), user.getGroupId());
            sendTextAsReply(chatId, fastestGG, msg.getMessageId());
            isFastestGGAvailable = false;
        }
    }

    /**
    * Determines if a message text is a gg
    *
    * @param text message text
    * @return true if message text is gg
    */
    public boolean isGG(String text) {
        return (Objects.equals(text, "GG") || Objects.equals(text, "gg") || Objects.equals(text, "Gg"));
    }

    /**
     * Parse the caption attached to the photo message that represents an activity
     * The format should be either activity-mm-dd-yyyy or activity-MMddyyyy
     *
     * @param msgCaption caption attached to the message
     * @param chatId the chatId for the telegram group
     * @return a map containing name, activity, day, month, year
     */
    public Map<String, String> parseMessageCaption(String msgCaption, Long chatId) {
        Map<String, String> parsedMessageContent = new HashMap<>();
        String[] splitCaption = msgCaption.split(DATE_SPLIT_FORMAT);

        if (splitCaption.length < 3) {
            sendText(chatId, "Invalid message format. Please follow the activity-mm-dd-yyyy or activity-MMddyyyy format.");
        }

        String activity = splitCaption[0];

        String month, day, year;

        if (splitCaption.length == 4) { //activity-mm-dd-yyyy
            month = splitCaption[1];
            day = splitCaption[2];
            year = splitCaption[3];
        } else { // activity-MMddyyyy
            if (splitCaption[1].length() != 8) {
                sendText(chatId, "Invalid message format. Please follow the activity-mm-dd-yyyy or activity-MMddyyyy format.");
            }

            DateTimeFormatter formatter = DateTimeFormatter.ofPattern(DATE_FORMAT, Locale.ENGLISH);
            LocalDate date = LocalDate.parse(splitCaption[1], formatter);
            String[] parsedDate = date.toString().split(DATE_SEPARATOR);

            year = parsedDate[0];
            month = parsedDate[1];
            day = parsedDate[2];
        }

        parsedMessageContent.put("activity", activity);
        parsedMessageContent.put("day", day);
        parsedMessageContent.put("month", month);
        parsedMessageContent.put("year", year);

        return parsedMessageContent;
    }

    /**
    * Gets the user associated with the message. If the user does not exist in
    * the database create a new user.
    *
    * @param msg the message
    * @return the user
    */
    public FitBoiUser getUserFromMessage(Message msg) {
        User user = msg.getFrom();

        Optional<FitBoiUser> fitBoiUser = userRepository.findById(user.getId());

        if (fitBoiUser.isPresent()) {
           return fitBoiUser.get();
        } else {
            FitBoiUser newFitBoiUser = new FitBoiUser(user.getId(), user.getFirstName(), msg.getChatId());
            FitBoiGg newFitBoiGg = new FitBoiGg(user.getId(), msg.getChatId(), String.valueOf(LocalDate.now().getYear()), 0);
            userRepository.save(newFitBoiUser);
            ggRepository.save(newFitBoiGg);
            return newFitBoiUser;
        }
    }

    /**
     * Gets the FitBoi record from the message caption. Save the record in the
     * database.
     *
     * @param captionContents parsed caption from the message
     */
    public void getRecordFromCaption(Long userId, Map<String, String> captionContents) {
        String activity = captionContents.get("activity");
        String month = captionContents.get("month");
        String day =  captionContents.get("day");
        String year = captionContents.get("year");

        if (year != null && year.length() == 2) {
            int yearInt = Integer.parseInt(year);
            int currentYear = Year.now().getValue();
            int currentYearLastTwoDigits = currentYear % 100;

            year = String.valueOf((currentYearLastTwoDigits < yearInt) ? currentYear - currentYearLastTwoDigits - 100 + yearInt : currentYear - currentYearLastTwoDigits + yearInt);
        }

        FitBoiRecord newFitBoiRecord = new FitBoiRecord(userId, activity, month, day, year);
        recordRepository.save(newFitBoiRecord);
    }

    /**
    * Using distinct names from the fit_boi_record table build an activity
    * counts message
    *
    * @return a message display activity counts
    */
    public String getActivityCountsMessage() {
        //TODO: Before returning activity counts make sure the fetched users belong in that
        //TODO: particular group.
        List<Long> userIds = recordRepository.findDistinctRecords();
        HashMap<String, Long> counts = new HashMap<>();

        ZonedDateTime nowInEST = ZonedDateTime.now(ZoneId.of("America/New_York"));
        System.out.println("ACTIVITY COUNTS MESSAGE - Time now in EST: " + nowInEST);
        String month = nowInEST.format(DateTimeFormatter.ofPattern("MM"));

        for (Long userId : userIds) {
            Long countOfRecords = recordRepository.countByUserIdWithCurrentYearAndMonth(userId, month);
            String name = userRepository.findById(userId).get().getName();

            counts.put(name, countOfRecords);
        }

        String content = counts.entrySet()
                .stream()
                .map(e -> e.getKey() + "=" + e.getValue())
                .collect(Collectors.joining(", "));

        return "Activity counts updated: " + content;
    }

    public String getFastestGgCountsMessage(Long chatId) {
        List<FitBoiUser> users = userRepository.findFitBoiUsersByGroupId(chatId);
        HashMap<String, Integer> counts = new HashMap<>();

        for (FitBoiUser user : users) {
            int ggCount = ggRepository.fetchFastGgCountByIdAndCurrentYear(user.getId());

            counts.put(user.getName(), ggCount);
        }

        String content = counts.entrySet()
                .stream()
                .sorted(Map.Entry.comparingByValue(Comparator.reverseOrder()))
                .map(e -> e.getKey() + "=" + e.getValue())
                .collect(Collectors.joining(", "));

        return "Fastest GGs \uD83D\uDE0E " + content;
    }
}
