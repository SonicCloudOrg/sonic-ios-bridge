package org.cloud.sonic.model;

public class ConnectStatus {
    public Device device;
    public Status status;

    public enum Status {
        ONLINE, DISCONNECT
    }
}
