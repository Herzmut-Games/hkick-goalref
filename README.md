# hkick-goalref
The `hkick-goalref` is part of the hkick project and is build to recognize any goal activity.

It is capable of monitoring two break beam sensors (we use the [adafruit IR sensors](https://www.adafruit.com/product/2167)). A goal is triggered via a interrupt if the ball passes the sensors installed in tunnel between goal and ball return area.

hkick is meant to be very modular so you could use a small device such a raspberry pi zero or esp8266 to only run the  `hkick-goalref` next to your kicker and have the rest of the software running somewhere else.

Currently `hkic-goalref` reports a goal via HTTP to `hkick-core`. This later will be switched to MQTT.
