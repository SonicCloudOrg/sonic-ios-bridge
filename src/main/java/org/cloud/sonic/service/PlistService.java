package org.cloud.sonic.service;

import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONObject;
import com.dd.plist.*;
import org.cloud.sonic.common.Tool;
import org.cloud.sonic.model.Device;
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

    private static int getSize() throws IOException {
        byte[] header = new byte[16];
        inputStream.read(header);
        ByteBuffer buffer = ByteBuffer.allocate(16);
        buffer.order(ByteOrder.LITTLE_ENDIAN);
        buffer.put(header);
        return buffer.getInt(0) - 16;
    }


    private static int getSize2(int readLen) throws IOException {
        byte[] header = new byte[readLen];
        inputStream.read(header);
        ByteBuffer buffer = ByteBuffer.allocate(readLen);
        buffer.order(ByteOrder.BIG_ENDIAN);
        buffer.put(header);
        return buffer.getInt(0);
    }

    public static NSDictionary getNsDictionary(int size) throws IOException, ParseException, ParserConfigurationException, SAXException, PropertyListFormatException {
        byte[] body = new byte[size];
        inputStream.read(body);
        NSObject parse = PropertyListParser.parse(body);
        return (NSDictionary) parse;
    }

    private static NSDictionary receiveMsg() throws IOException, PropertyListFormatException, ParseException, ParserConfigurationException, SAXException {
        byte[] res = new byte[4];
        inputStream.read(res);
        ByteBuffer readB = ByteBuffer.allocate(res.length);
        readB.order(ByteOrder.BIG_ENDIAN);
        readB.put(res);
        int aInt = readB.getInt(0);
        byte[] body = new byte[aInt];
        inputStream.read(body);
        NSObject parse = PropertyListParser.parse(body);
        System.out.println(parse.toXMLPropertyList());
        return (NSDictionary) parse;
    }

    private static ByteBuffer buildByteMsg(byte[] bytes) throws IOException {
        int len = (16 + bytes.length);
        int version = 1;
        int request = 8;
        int tag = 1;
        ByteBuffer buffer = ByteBuffer.allocate(len);
        buffer.order(ByteOrder.LITTLE_ENDIAN);
        buffer.putInt(0, len);
        buffer.putInt(4, version);
        buffer.putInt(8, request);
        buffer.putInt(12, tag);
        int i = 16;
        for (byte aByte : bytes) {
            buffer.put(i++, aByte);
        }
        return buffer;
    }

    private static NSDictionary sendAndReceiveMsg() throws IOException, PropertyListFormatException, ParseException, ParserConfigurationException, SAXException {
        NSDictionary root = new NSDictionary();
        root.put("MessageType", "ListDevices");
        root.put("ClientVersionString", "libusbmuxd 1.1.0");
        root.put("ProgName", "sonic-ios-bridge");
        root.put("kLibUSBMuxVersion", "3");
        String s = root.toXMLPropertyList();
        outputStream.write(buildByteMsg(s.getBytes(Charset.forName("UTF-8"))).array());
        int size = getSize();
        NSDictionary dico = getNsDictionary(size);
        return dico;
    }

    private static NSDictionary connectDevice(int id, int port) throws IOException, PropertyListFormatException, ParseException, ParserConfigurationException, SAXException {
        NSDictionary root = new NSDictionary();
        root.put("MessageType", "Connect");
//        root.put("ClientVersionString", "sonic-ios-bridge");
        root.put("ProgName", "sonic-ios-bridge");
        root.put("DeviceID", new NSNumber(id));
        root.put("PortNumber", new NSNumber(swapPortNumber(port)));
        String s = root.toXMLPropertyList();
        outputStream.write(buildByteMsg(s.getBytes(Charset.forName("UTF-8"))).array());
        int size = getSize();
        NSDictionary dico = getNsDictionary(size);
        return dico;
    }

    private static NSDictionary sendAndReceiveMsg2() throws IOException, PropertyListFormatException, ParseException, ParserConfigurationException, SAXException {
        NSDictionary root = new NSDictionary();
        root.put("Request", "GetValue");
//        root.put("Key", "ProductVersion");
        root.put("Domain", "com.apple.mobile.iTunes");
        root.put("Label", "sonic-ios-bridge");
//        root.put("ProgName", "sonic-ios-bridge");
//        root.put("kLibUSBMuxVersion", "3");
        String s = root.toXMLPropertyList();
        byte[] bytes = s.getBytes(Charset.forName("UTF-8"));
        ByteBuffer buffer = ByteBuffer.allocate(
                4);
        buffer.order(ByteOrder.BIG_ENDIAN);
        buffer.putInt(bytes.length);
        outputStream.write(buffer.array());
        buffer = ByteBuffer.allocate(
                bytes.length);
        buffer.order(ByteOrder.BIG_ENDIAN);
        buffer.put(bytes);
        outputStream.write(bytes);
//        sendMsg(s.getBytes(Charset.forName("UTF-8")), -1);
        return receiveMsg();
    }

//    private static NSDictionary sendAndReceiveMsg2() throws IOException, PropertyListFormatException, ParseException, ParserConfigurationException, SAXException {
//        NSDictionary root = new NSDictionary();
//        root.put("Request", "QueryType");
////        root.put("Key", "ProductVersion");
////        root.put("Domain", "com.apple.mobile.iTunes");
////        root.put("Label", "sonic-ios-bridge");
////        root.put("ProgName", "sonic-ios-bridge");
////        root.put("kLibUSBMuxVersion", "3");
//        String s = root.toXMLPropertyList();
//        byte[] bytes = s.getBytes(Charset.forName("UTF-8"));
//        ByteBuffer buffer = ByteBuffer.allocate(
//                4);
//        buffer.order(ByteOrder.BIG_ENDIAN);
//        buffer.putInt(bytes.length);
//        outputStream.write(buffer.array());
//        buffer = ByteBuffer.allocate(
//                bytes.length);
//        buffer.order(ByteOrder.BIG_ENDIAN);
//        buffer.put(bytes);
//        outputStream.write(bytes);
////        sendMsg(s.getBytes(Charset.forName("UTF-8")), -1);
//        return receiveMsg();
//    }

    protected static int swapPortNumber(int port) {
        return ((port << 8) & 0xFF00) | (port >> 8);
    }

    public static void main(String[] args) throws IOException, PropertyListFormatException, ParseException, ParserConfigurationException, SAXException {
        connect();
        NSArray a = (NSArray) sendAndReceiveMsg().get("DeviceList");
        NSDictionary b = (NSDictionary) a.objectAtIndex(0);
        Device c = b.get("Properties").toJavaObject(Device.class);
//        System.out.println(b.get("Properties").toXMLPropertyList());
//        System.out.println(c);
//        connect();
//        isFirstPack = true;
        System.out.println(connectDevice(c.getDeviceID(), 62078).toXMLPropertyList());
//        sendAndReceiveMsg2()
//        isFirstPack = true;
        sendAndReceiveMsg2().toXMLPropertyList();
//

    }
}
