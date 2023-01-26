package com.alcamech.fitboisbot;

public enum Commands {
    HELP("/help"),
    FASTGG("/fastgg");

    private final String text;

    Commands(final String text) {
        this.text = text;
    }

    @Override
    public String toString() {
        return text;
    }
}
