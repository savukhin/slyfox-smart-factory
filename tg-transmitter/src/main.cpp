#include <Arduino.h>
#include "defines.hpp"
#include <FastBot.h>
#include <set>

#include "imu_arduino.hpp"
#include "mic_arduino.hpp"
#include "imu_stub.hpp"
#include "mic_stub.hpp"
#include "buffer.hpp"

FastBot bot(BOT_TOKEN);
unsigned long prev_connect = 0;
unsigned long time_notconnected = 0;
unsigned long time_connected = 0;
uint32_t startUnix;
uint32_t startUnixMillis;
std::set<String> sendDataChatIDs;

// ImuArduino imu;
// MicArduino mic(MIC_PIN);
ImuStub imu;
MicStub mic;
Buffer buffer(100);

void connectWiFi()
{
	delay(2000);
	WiFi.begin(WIFI_SSID, WIFI_PASS);
	prev_connect = millis();
	while (WiFi.status() != WL_CONNECTED)
	{
		for (unsigned short x = 0; x <= 2; x++)
		{
			digitalWrite(2, HIGH);
			delay(250);
			digitalWrite(2, LOW);
			delay(250);
		}
		if ((millis() - prev_connect) > 5000)
			return;
	}
	Serial.println("Connected");
}

void stat_wifi()
{
	if (WiFi.status() != WL_CONNECTED)
	{
		digitalWrite(LED_BUILTIN, HIGH);
		delay(50);
		digitalWrite(LED_BUILTIN, LOW);
		delay(50);
		time_notconnected = millis();
	}
	if (WiFi.status() == WL_CONNECTED)
		time_connected = millis();

	if (WiFi.status() != WL_CONNECTED && ((time_notconnected - time_connected) >= 60000) && (millis() - prev_connect) > 15000)
		connectWiFi();
}


void newMsg(FB_msg &msg)
{
	Serial.println("got msg" + msg.text);
	if (msg.unix < startUnix)
		return;

	if (msg.text == "/hardreset")
	{
		bot.sendMessage("ПЕРЕЗАГРУЖАЮСЬ...", msg.chatID);
		delay(1000);
		bot.tickManual();
		ESP.restart();
	}

	if (msg.text == "/data")
	{
		sendDataChatIDs.insert(msg.chatID);
		bot.sendMessage("Теперь сообщения с данными будут отправляться в этот чат", msg.chatID);
	}
}

void setup()
{
	Serial.begin(115200);
	connectWiFi();
	delay(1000);
	bot.setChatID("");
	bot.skipUpdates();
	startUnix = bot.getUnix();
	startUnixMillis = millis();
	bot.attach(newMsg);
}

String serializeCsv(BufferItems &items) {
	String s = "timestamp,accX(m/s^2),accY,accZ,gyroX(rad/s),gyroY,gyroZ,mic\n";
	for (int i = 0; i < items.size; i++) {
		BufferItem item = items.items_[i];
		s += String(item.timestamp) + ",";
		s += String(item.imu.acc.x) + ",";
		s += String(item.imu.acc.y) + ",";
		s += String(item.imu.acc.z) + ",";
		s += String(item.imu.gyro.x) + ",";
		s += String(item.imu.gyro.y) + ",";
		s += String(item.imu.gyro.z) + ",";
		s += String(item.mic.value);
		s += "\n";
	}
	return s;
}

bool needSendData() {
	// Serial.printf("needSendData");
	// Serial.printf("sendDataChatIDs.empty %d\n", sendDataChatIDs.empty());
	return !sendDataChatIDs.empty();
	// return false;
}

String getChatIDs() {
	String result;
	for (auto id : sendDataChatIDs) {
		result = result + "," + id;
	}
	return result.substring(1);
}

void sendItems(BufferItems &items) {
	// Serial.println("send items");
	if (!needSendData())
		return;

	String str = serializeCsv(items);
	uint8_t *s = (uint8_t*)str.c_str();
	uint8_t status = bot.sendFile(s, str.length(), FB_DOC, String(millis()) + ".csv", getChatIDs());

	if (status != 1) {
		Serial.printf("Got fail status %d", status);
	}
}

void loop()
{
	bot.tick();
	// Serial.printf("loop %d! %d\n", buffer.size(), WiFi.status());
	ImuData imuData = imu.GetData();
	MicData micData = mic.GetData();
	BufferItems flushedItems = buffer.insert(imuData, micData);
	// Serial.printf("if %d\t", flushedItems.size);
	if (flushedItems.size == 0)
		return;
	
	sendItems(flushedItems);

}
