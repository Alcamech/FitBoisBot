package com.alcamech.fitboisbot;

import com.alcamech.fitboisbot.model.FitBoiRecord;
import com.alcamech.fitboisbot.model.FitBoiUser;
import com.alcamech.fitboisbot.respository.FitBoisRepository;
import com.alcamech.fitboisbot.respository.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.telegram.telegrambots.bots.TelegramLongPollingBot;
import org.telegram.telegrambots.meta.api.methods.send.SendMessage;
import org.telegram.telegrambots.meta.api.objects.Message;
import org.telegram.telegrambots.meta.api.objects.Update;
import org.telegram.telegrambots.meta.api.objects.User;
import org.telegram.telegrambots.meta.exceptions.TelegramApiException;

import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.util.*;
import java.util.stream.Collectors;

@Component
public class FitBoisBot extends TelegramLongPollingBot {

    @Autowired
    FitBoisRepository fitBoisRepository;
    @Autowired
    UserRepository userRepository;
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

    @Override
    public void onUpdateReceived(Update update) {
        //TODO: clean this method up, break it up, add better error handling
        Message msg;
        if (update.getEditedMessage() != null) {
            msg = update.getEditedMessage();
        } else {
            msg = update.getMessage();

            if (msg.hasPhoto()) { // photo post is potentially an activity
                lastActivityPostUserId = msg.getFrom().getId();
                isFastestGGAvailable = true;
            }
        }

        FitBoiUser user = getUserFromMessage(msg);
        Long chatId = msg.getChat().getId();

        isFastestGG(update, msg, chatId, user);

        if (msg.hasPhoto() && msg.getCaption() != null) {
            String msgCaption = msg.getCaption();

            Map<String, String> captionContents = parseMessageCaption(msgCaption, chatId);
            getRecordFromCaption(captionContents);

            String totalActivitiesMessage = getActivityCountsMessage();
            sendText(chatId, totalActivitiesMessage);
        }

        if (msg.isCommand()) {
            if (msg.getText().equals(Commands.HELP.toString())) {
                onHelp(chatId);
            }

            if (msg.getText().equals(Commands.FASTGG.toString())) {
                onFastGG(chatId);
            }
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
    * Response for the /fastgg command
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
        // fastest GG is not available;
        if (!isFastestGGAvailable) {
            return;
        }

        // update is an edit from activity poster
        if (update.getEditedMessage() != null && update.getEditedMessage().getFrom().getId().equals(lastActivityPostUserId)) {
            return;
        }

        // message is from activity poster
        if (msg.getFrom().getId().equals(lastActivityPostUserId)) {
           return;
        }

        if (isGG(msg.getText())) {
            userRepository.updateGgCount(user.getId(), user.getGroupId());
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
    * The format should be either name-activity-mm-dd-yyyy or name-activity-MMddyyyy
    *
    * @param msgCaption caption attached to the message
    * @param chatId the chatId for the telegram group
    * @return a map containing name, activity, day, month, year
    */
    public Map<String, String> parseMessageCaption(String msgCaption, Long chatId) {
        Map<String, String> parsedMessageContent = new HashMap<>();

        String[] splitCaption = msgCaption.split("-");

        String name = "", activity = "";

        try {
            name = splitCaption[0];
            activity = splitCaption[1];
        } catch (Exception e) {
            sendText(chatId, "Something went wrong. Check your message formatting.");
        }

        String month = "", day = "", year = "";

        if (splitCaption.length == 5) { //name-activity-mm-dd-yyyy
            month = splitCaption[2];
            day = splitCaption[3];
            year = splitCaption[4];
        } else  { // name-activity-MMddyyyy
            DateTimeFormatter formatter = DateTimeFormatter.ofPattern("MMddyyyy", Locale.ENGLISH);
            LocalDate date = LocalDate.parse(splitCaption[2], formatter);
            String[] parsedDate = date.toString().split("-");

            month = parsedDate[1];
            day = parsedDate[2];
            year = parsedDate[0];
        }

        parsedMessageContent.put("name", name);
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
            userRepository.save(newFitBoiUser);
            userRepository.initializeGgCount(user.getId(), msg.getChatId());

            return newFitBoiUser;
        }
    }

    /**
    * Gets the FitBoi record from the message caption. Save the record in the
    * database.
    *
    * @param captionContents parsed caption from the message
    * @return the FitBoi record
    */
    public FitBoiRecord getRecordFromCaption(Map<String, String> captionContents) {
        String name = captionContents.get("name");
        String activity = captionContents.get("activity");
        String month = captionContents.get("month");
        String day =  captionContents.get("day");
        String year = captionContents.get("year");

        FitBoiRecord newFitBoiRecord = new FitBoiRecord(name, activity, month, day, year);
        fitBoisRepository.save(newFitBoiRecord);

        return newFitBoiRecord;
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
        List<String> names = fitBoisRepository.findDistinctName();
        HashMap<String, Long> counts = new HashMap<>();

        for (String retrievedName : names) {
            Long countOfRecords = fitBoisRepository.countByName(retrievedName);

            counts.put(retrievedName, countOfRecords);
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
            int ggCount = userRepository.fetchFastGgCountById(user.getId());

            counts.put(user.getName(), ggCount);
        }

        String content = counts.entrySet()
                .stream()
                .map(e -> e.getKey() + "=" + e.getValue())
                .collect(Collectors.joining(", "));

        return "Fastest GGs \uD83D\uDE0E " + content;
    }
}