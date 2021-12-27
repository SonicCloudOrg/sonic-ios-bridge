package org.cloud.sonic.service;

import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONObject;
import com.dd.plist.NSDictionary;
import com.dd.plist.NSObject;
import com.dd.plist.PropertyListFormatException;
import com.dd.plist.PropertyListParser;
import org.cloud.sonic.common.Tool;
import org.newsclub.net.unix.AFUNIXSocket;
import org.newsclub.net.unix.AFUNIXSocketAddress;
import org.xml.sax.SAXException;

import javax.xml.parsers.ParserConfigurationException;
import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.net.Socket;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.charset.Charset;
import java.text.ParseException;

public class PlistService {
    private static String system = System.getProperty("os.name").toLowerCase();
    private static String connectAddress;
    private static Socket socket;
    private static InputStream inputStream;
    private static OutputStream outputStream;
    private static Boolean isFirstPack = true;

    private static void getAddress() {
        if (system.contains("win")) {
            connectAddress = "127.0.0.1:27015";
        } else if (system.contains("linux") || system.contains("mac")) {
            connectAddress = "/var/run/usbmuxd";
        } else {
            throw new RuntimeException(String.format("This os is not support: %s , " +
                    "you can commit issue to https://github.com/SonicCloudOrg/sonic-ios-bridge/issues", system));
        }
    }

    private static void connect() throws IOException {
        getAddress();
        if (connectAddress.contains(":")) {
            socket = new Socket();
            socket.connect(new InetSocketAddress(
                    connectAddress.substring(0, connectAddress.indexOf(":")),
                    Integer.parseInt(connectAddress.substring(connectAddress.indexOf(":") + 1))));
            inputStream = socket.getInputStream();
            outputStream = socket.getOutputStream();
        } else if (system.contains("linux") || system.contains("mac")) {
            socket = AFUNIXSocket.newInstance();
            socket.connect(new AFUNIXSocketAddress(new File("/var/run/usbmuxd")));
            inputStream = socket.getInputStream();
            outputStream = socket.getOutputStream();
        } else {
            throw new RuntimeException(String.format("This file is not exist: %s , " +
                    "please check again"));
        }
    }

//    private byte[] receiveAllMsg(int size) {
//        byte[] buffer = new byte[0];
//        while (buffer.length < size) {
//            byte[] chunk = new byte[size - buffer.length];
//            int realLen = 0;
//            try {
//                realLen = inputStream.read(chunk);
//            } catch (IOException e) {
//                e.printStackTrace();
//            }
//            if (chunk.length != realLen && realLen >= 0) {
//                buffer = Tool.subByteArray(chunk, 0, realLen);
//            }
//        }
//        return buffer;
//    }
//
//    private void sendAllMsg(byte[] bytes) throws IOException {
//        outputStream.write(bytes);
//        outputStream.flush();
//    }

    private static int getSize(int readLen) throws IOException {
        byte[] header = new byte[readLen];
        inputStream.read(header);
        ByteBuffer buffer = ByteBuffer.allocate(readLen);
        buffer.order(ByteOrder.LITTLE_ENDIAN);
        buffer.put(header);
        return buffer.getInt(0);
    }

    public static NSDictionary getNsDictionary(int size) throws IOException, ParseException, ParserConfigurationException, SAXException, PropertyListFormatException {
        byte[] body = new byte[size];
        inputStream.read(body);
        NSObject parse = PropertyListParser.parse(body);
        return (NSDictionary) parse;
    }

    private static NSDictionary receiveMsg(int headerLen) throws IOException, PropertyListFormatException, ParseException, ParserConfigurationException, SAXException {
        int length;
        if (isFirstPack || headerLen == 16) {
            length = getSize(16);
            length -= 16;
            isFirstPack = false;
        } else {
            length = getSize(4);
        }
        return getNsDictionary(length);
    }

    private static void sendMsg(byte[] bytes, int msgType) throws IOException {
        int version = 1;
        int request = msgType != -1 ? msgType : 8;
        int tag = 1;
        ByteBuffer buffer;
        if (isFirstPack) {
            int len = (16 + bytes.length);
            buffer = ByteBuffer.allocate(len);
            buffer.order(ByteOrder.LITTLE_ENDIAN);
            buffer.putInt(0, len);
            buffer.putInt(4, version);
            buffer.putInt(8, request);
            buffer.putInt(12, tag);
            int i = 16;
            for (byte aByte : bytes) {
                buffer.put(i++, aByte);
            }
        } else {
            buffer = ByteBuffer.allocate(bytes.length);
            buffer.order(ByteOrder.BIG_ENDIAN);
            buffer.putInt(0, bytes.length);
            int i = 0;
            for (byte aByte : bytes) {
                buffer.put(i++, aByte);
            }
        }
        outputStream.write(buffer.array());
    }

    private static NSDictionary sendAndReceiveMsg() throws IOException, PropertyListFormatException, ParseException, ParserConfigurationException, SAXException {
        NSDictionary root = new NSDictionary();
        root.put("MessageType", "ListDevices");
        root.put("ClientVersionString", "sonic-ios-bridge");
        root.put("ProgName", "sonic-ios-bridge");
        root.put("kLibUSBMuxVersion", "3");
        String s = root.toXMLPropertyList();
        sendMsg(s.getBytes(Charset.forName("UTF-8")),-1);
        return receiveMsg(4);
    }

    public static void main(String[] args) throws IOException, PropertyListFormatException, ParseException, ParserConfigurationException, SAXException {
        connect();
        String[] a = sendAndReceiveMsg().allKeys();
        for(String b:a){
            System.out.println(b);
        }
    }
}
