# Mattermost Antivirus Plugin (Beta) ![CircleCI branch](https://img.shields.io/circleci/project/github/mattermost/mattermost-plugin-antivirus/master.svg)

This plugin allows the forwarding of uploaded files to an antivirus scanning application and prevents the upload from completing if there is a virus detected in the file. Use it to prevent users from inadvertently spreading malware or viruses via your Mattermost instance. 

Currently the plugin supports [ClamAV anti-virus software](https://www.clamav.net/) across browser, Desktop Apps and the Mobile Apps. ClamAV is an open source (GPL) anti-virus engine used in a variety of situations including email scanning, web scanning, and end point security. It provides a number of utilities including a flexible and scalable multi-threaded daemon and an advanced tool for automatic database updates. A ClamAV server can be easily provisioned as a docker container that runs alongside Mattermost. 

**Requirements:**

- Mattermost Server Version: 5.2+
- ClamAV Server access

## Installation

1. Go to the [releases page of this Github repository](https://github.com/mattermost/mattermost-plugin-antivirus/releases) and download the latest release for your Mattermost server.
2. In the Mattermost System Console under **System Console > Plugins > Management** upload the file to install the plugin. To learn more about how to upload a plugin, [see the documentation](https://docs.mattermost.com/administration/plugins.html#plugin-uploads).
3. Install ClamAV (clamd) for virus scanning. One easy option is to provision a ClamAV container with docker by running the following command.  Assuming you have already installed Docker, this will download and install the latest version of ClamAV and setup a server with an open port at 3310:

   ```
   docker run -d -p 3310:3310 mkodockx/docker-clamav
   ```

4. Once clamd server is running, configure the plugin to make requests to your clamd instance. Go to **System Console > Plugins > Antivirus** and configure **Clamav Host and Port** to point at your clamd instance.  
5. Activate the plugin at **System Console > Plugins > Management** and ensure that it starts with no errors


## Testing

To test your configuration is correct, create an [EICAR test file](https://2016.eicar.org/86-0-Intended-use.html) (copy the text from that webpage into a text file editor and save it) and upload it. The file should be rejected as below:

![Screenshot of Anti-virus in action](/2019-07-26_13-56-13.png)

Upload a regular file to ensure they are processed OK and posted to the channel.

If there is an error with your setup - check your ClamAV server setup and communication:

![Screenshot of Anti-virus plugin showing a server error](/2019-07-26_11-52-33.png)




