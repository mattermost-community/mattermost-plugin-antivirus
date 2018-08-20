# Mattermost Antivirus Plugin (Beta)

This plugin allows the forwarding of uploaded files to an antivirus application. Use it to scan for viruses before uploading a file to Mattermost.

Currently the plugin supports [ClamAV anti-virus software](https://www.clamav.net/) across browser, Desktop Apps and the Mobile Apps.

**Supported Mattermost Server Versions: 5.2+**

## Installation

1. Go to the [releases page of this Github repository](https://github.com/mattermost/mattermost-plugin-antivirus/releases) and download the latest release for your Mattermost server.
2. Upload this file in the Mattermost System Console under **System Console > Plugins > Management** to install the plugin. To learn more about how to upload a plugin, [see the documentation](https://docs.mattermost.com/administration/plugins.html#plugin-uploads).
3. Install ClamAV (clamd) for virus scanning. One option is to install with docker:

   ```
   docker run -d -p 3310:3310 mkodockx/docker-clamav
   ```

4. Once clamd is running, configure the plugin to make requests to your clamd instance. Go to **System Console > Plugins > Antivirus** and configure **Clamav Host and Port** to point at your clamd instance.
5. Activate the plugin at **System Console > Plugins > Management**.

You're all set! All file uploads on the system are now scanned for viruses by ClamAV.

To test your configuration is correct, create an [EICAR test file](https://www.eicar.org/86-0-Intended-use.html) and upload it. The file should be rejected.
