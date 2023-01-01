package com.alcamech.fitboisbot;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.telegram.telegrambots.bots.TelegramLongPollingBot;
import org.telegram.telegrambots.meta.api.methods.send.SendMessage;
import org.telegram.telegrambots.meta.api.objects.Message;
import org.telegram.telegrambots.meta.api.objects.Update;
import org.telegram.telegrambots.meta.exceptions.TelegramApiException;

import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Locale;
import java.util.stream.Collectors;

@Component
public class FitBoisBot extends TelegramLongPollingBot {

    @Autowired
    FitBoisRepository fitBoisRepository;

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

    @Override
    public void onUpdateReceived(Update update) {
        //TODO: Break this method up and add better error handling
        Message msg;

        if (update.getEditedMessage() != null) {
            msg = update.getEditedMessage();
        } else {
            msg = update.getMessage();
        }

        Long chatId = msg.getChat().getId();

        if (msg.hasPhoto() && msg.getCaption() != null) {
            String activityContent = msg.getCaption();

            String[] parsedActivityContent = activityContent.split("-");

            String name = "", activity = "";

            try {
                name = parsedActivityContent[0];
                activity = parsedActivityContent[1];
            } catch (Exception e) {
                sendText(chatId, "Something went wrong. Check your message formatting.");
            }

            String month = "", day = "", year = "";

            if (parsedActivityContent.length == 5) { //name-activity-mm-dd-yyyy
                month = parsedActivityContent[2];
                day = parsedActivityContent[3];
                year = parsedActivityContent[4];
            } else  { // name-activity-MMddyyyy
                DateTimeFormatter formatter = DateTimeFormatter.ofPattern("MMddyyyy", Locale.ENGLISH);
                LocalDate date = LocalDate.parse(parsedActivityContent[2], formatter);
                String[] parsedDate = date.toString().split("-");

                month = parsedDate[1];
                day = parsedDate[2];
                year = parsedDate[0];
            }

            FitBoisController controller = new FitBoisController(fitBoisRepository);
            controller.addNewRecord(name, activity, month, day, year);

            List<String> names = controller.getFitBois();
            HashMap<String, Long> counts = new HashMap<>();

            for (String retrievedName : names) {
                Long countOfRecords = controller.getCountByName(retrievedName);

                counts.put(retrievedName, countOfRecords);
            }

            String content = counts.entrySet()
                    .stream()
                    .map(e -> e.getKey() + "=" + e.getValue())
                    .collect(Collectors.joining(", "));

            String totalActivitiesMessage = "Activity counts updated: " + content;
            sendText(chatId, totalActivitiesMessage);
        }

        if (msg.isCommand()) {
            if (msg.getText().equals(Commands.HELP.toString())) {
                onHelp(chatId);
            }
        }
    }

    /**
    * Response for the /help command
    *
    * @param chatId chat id to send response to
    */
    public void onHelp(Long chatId ) {
        String message = "Nothing to see here... yet";
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
}