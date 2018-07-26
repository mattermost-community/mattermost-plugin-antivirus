# Mattermost Antivirus Plugin

This plugin allows the forwarding of uploaded files to an antivirus application. Currently the plugin supports clamav.

Requires Mattermost 5.2 or above.

## Installation

Go to the [releases page of this Github repository](https://github.com/mattermost/mattermost-plugin-antivirus/releases) and download the latest release. You can upload this file in the Mattermost system console under **System Console > Plugins > Management** to install the plugin. For more details on uploading and installing plugins, [see the full documentation.](https://docs.mattermost.com/administration/plugins.html#plugin-uploads).

You will need a clamd installation to do the actual scanning. An easy way to get one is with docker:

```
docker run -d -p 3310:3310 mkodockx/docker-clamav
```

Once you have clamd running you will need to configure the plugin to make requests to your clamd instance.  Go to **System Console > Plugins > Antivirus** and configure "Clamav Host and Port" to point at your clamd instance.

You can then activate the plugin **System Console > Plugins > Management** and all your uploads will be scanned.

You can test that you have set everything up correctly by making an [EICAR test file](https://www.eicar.org/86-0-Intended-use.html) and uploading it. It should be rejected.
