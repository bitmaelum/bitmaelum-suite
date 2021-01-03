# BitMaelum client

A docker image with the BitMaelum client software that allows you to easily test drive the BitMaelum mailing network. In order to work, you need to mount a local 
directory to the `/bitmaelum` directory:

      $ docker run -ti $PWD:/bitmaelum bitmaelum/client:latest account list


### Environment settings

| Settings                | Description                            |
|-------------------------|----------------------------------------|
| BITMAELUM_CLIENT_CONFIG | Path to custom client configuration    |
