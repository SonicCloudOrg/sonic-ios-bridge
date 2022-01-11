package org.cloud.sonic.model;

public class Device {
	public String udId;
	public String serialNumber;
	public Integer productId;
	public Integer locationId;
	public Integer deviceId;
	public String connectionType;

	public String getUDID() {
		return udId;
	}

	public void setUDID(String udId) {
		this.udId = udId;
	}

	public String getSerialNumber() {
		return serialNumber;
	}

	public void setSerialNumber(String serialNumber) {
		this.serialNumber = serialNumber;
	}

	public Integer getProductID() {
		return productId;
	}

	public void setProductID(Integer productId) {
		this.productId = productId;
	}

	public Integer getLocationID() {
		return locationId;
	}

	public void setLocationID(Integer locationId) {
		this.locationId = locationId;
	}

	public Integer getDeviceID() {
		return deviceId;
	}

	public void setDeviceID(Integer deviceId) {
		this.deviceId = deviceId;
	}

	public String getConnectionType() {
		return connectionType;
	}

	public void setConnectionType(String connectionType) {
		this.connectionType = connectionType;
	}

	@Override
	public String toString() {
		return "Device{" +
				"udId='" + udId + '\'' +
				", serialNumber='" + serialNumber + '\'' +
				", productId=" + productId +
				", locationId=" + locationId +
				", deviceId=" + deviceId +
				", connectionType='" + connectionType + '\'' +
				'}';
	}
}
