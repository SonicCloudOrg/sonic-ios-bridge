package org.cloud.sonic.common;

public class Tool {
    public static int[] unpack(final byte[] byte_array) {
        final int[] integerReadings = new int[byte_array.length / 2];
        for(int counter = 0, integerCounter = 0; counter < byte_array.length;) {
            integerReadings[integerCounter] = convertTwoBytesToInteger(byte_array[counter], byte_array[counter + 1]);
            counter += 2;
            integerCounter++;
        }
        return integerReadings;
    }

    private static int convertTwoBytesToInteger(final byte byte1, final byte byte2) {
        final int unsignedInteger1 = getUnsignedInteger(byte1);
        final int unsignedInteger2 = getUnsignedInteger(byte2);
        return unsignedInteger1 * 256 + unsignedInteger2;
    }

    private static int getUnsignedInteger(final byte b) {
        int unsignedInteger = b;
        if(b < 0) {
            unsignedInteger = b + 256;
        }
        return unsignedInteger;
    }

    public static byte[] subByteArray(byte[] byte1, int start, int end) {
        byte[] byte2;
        byte2 = new byte[end - start];
        System.arraycopy(byte1, start, byte2, 0, end - start);
        return byte2;
    }
}
