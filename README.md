# Disclaimer

**This repository is community supported and not maintained by Mattermost. Mattermost disclaims liability for integrations, including Third Party Integrations and Mattermost Integrations. Integrations may be modified or discontinued at any time.**

# Mattermost Antivirus Plugin

![CircleCI branch](https://img.shields.io/circleci/project/github/mattermost/mattermost-plugin-antivirus/master.svg)
[![Code Coverage](https://img.shields.io/codecov/c/github/mattermost/mattermost-plugin-antivirus/master)](https://codecov.io/gh/mattermost/mattermost-plugin-antivirus)
[![Release](https://img.shields.io/github/v/release/mattermost/mattermost-plugin-antivirus)](https://github.com/mattermost/mattermost-plugin-antivirus/releases/latest)
[![HW](https://img.shields.io/github/issues/mattermost/mattermost-plugin-antivirus/Up%20For%20Grabs?color=dark%20green&label=Help%20Wanted)](https://github.com/mattermost/mattermost-plugin-antivirus/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22Up+For+Grabs%22+label%3A%22Help+Wanted%22)

**Maintainer:** [@iomodo](https://github.com/iomodo)
**Co-Maintainer:** [@hanzei](https://github.com/hanzei)

This plugin allows the forwarding of uploaded files to an antivirus scanning application and prevents the upload from completing if there is a virus detected in the file. Use it to prevent users from inadvertently spreading malware or viruses via your Mattermost instance. 

Currently the plugin supports [ClamAV anti-virus software](https://www.clamav.net/) across browser, Desktop Apps and the Mobile Apps. ClamAV is an open source (GPL) anti-virus engine used in a variety of situations including email scanning, web scanning, and end point security. It provides a number of utilities including a tool for automatic database updates. A ClamAV server can be easily provisioned as a Docker container that runs alongside Mattermost. 

**Requirements:**

- Mattermost Server Version: 5.2+
- ClamAV Server access

## Installation

1. Go to the [releases page of this Github repository](https://github.com/mattermost/mattermost-plugin-antivirus/releases) and download the latest release for your Mattermost server.
2. In the Mattermost System Console under **System Console > Plugins > Plugin Management** upload the file to install the plugin. To learn more about how to upload a plugin, [see the documentation](https://docs.mattermost.com/administration/plugins.html#plugin-uploads).

3. Install ClamAV (clamd) for virus scanning. One easy option is to provision a ClamAV container with Docker by running the following command.  Assuming you have already installed Docker, this will download and install the latest version of ClamAV and set up a server with an open port at 3310. ClamAV by default accepts 100MB files, you can change this in the [`clamd.conf`](https://github.com/Cisco-Talos/clamav/blob/main/etc/clamd.conf.sample). Visit the [ClamAV Documentation](https://docs.clamav.net/manual/Installing/Docker.html) for further configuration. Mattermost's MaxFileSize default value is subject to change. To ensure that the correct value is set, verify your value at the following link: [Maximum File Size](https://docs.mattermost.com/administration/config-settings.html#maximum-file-size).  

   If your Mattermost's MaxFileSize is ≤ 100MB
      ```
      docker run -d --restart unless-stopped -p 3310:3310 clamav/clamav:latest 
      ```
   If it is > 100MB
      ```
      docker run -d --restart unless-stopped --mount type=bind,source=/full/path/to/clamav/,target=/etc/clamav -p 3310:3310 clamav/clamav:latest 
      ```
   `/full/path/to/clamav/clamd.conf`
   ```conf
   ...
   # Files larger than this limit won't be scanned. Affects the input file itself
   # as well as files contained inside it (when the input file is an archive, a
   # document or some other kind of container).
   # Value of 0 disables the limit.
   # Note: disabling this limit or setting it too high may result in severe damage
   # to the system.
   # Technical design limitations prevent ClamAV from scanning files greater than
   # 2 GB at this time.
   # Default: 100M
   MaxFileSize 200M # Match your filezize limit in Mattermost
   ...
   ```

4. Once clamd server is running, configure the plugin in Mattermost to make requests to your clamd instance by going to **System Console > Plugins > Antivirus**. Configure **Clamav Host and Port** to point at your clamd instance, and optionally configure a **Scan timeout in seconds** to set how long it takes before the virus scan times out.  
5. Activate the plugin at **System Console > Plugins > Management** and ensure it starts with no errors.

## Testing

To test your configuration is correct, download an [EICAR test file](https://www.eicar.org/download-anti-malware-testfile/) and upload it. The file should be rejected as below:

![Screenshot of Anti-virus in action](/2019-07-26_13-56-13.png)

Upload a regular file to ensure it is processed successfully and posted to the channel.

If there is an error with your setup - check your ClamAV server setup and communication:

![Screenshot of Anti-virus plugin showing a server error](/2019-07-26_11-52-33.png)

## Development

This plugin contains both a server and web app portion. Read our documentation about the [Developer Workflow](https://developers.mattermost.com/integrate/plugins/developer-workflow/) and [Developer Setup](https://developers.mattermost.com/integrate/plugins/developer-setup/) for more information about developing and extending plugins.

### Releasing new versions

The version of a plugin is determined at compile time, automatically populating a `version` field in the [plugin manifest](plugin.json):
* If the current commit matches a tag, the version will match after stripping any leading `v`, e.g. `1.3.1`.
* Otherwise, the version will combine the nearest tag with `git rev-parse --short HEAD`, e.g. `1.3.1+d06e53e1`.
* If there is no version tag, an empty version will be combined with the short hash, e.g. `0.0.0+76081421`.

To disable this behaviour, manually populate and maintain the `version` field.
