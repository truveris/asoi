# asoi - AWS Status over IRC

This program monitors the AWS status feeds and updates IRC channels about it.

## Usage
Without any parameter, asoi uses `/etc/asoi.json` as configuration file:

	asoi [-c config]

## Example configuration
A configuration file is required to run asoi, this is an example:

	{
		"IRCServer": "irc.oftc.net:6667",
		"IRCNickname": "aws",
		"Channels": ["#dev"],
		"Topics": [
			"ec2-us-east-1",
			"sqs-us-east-1",
			"s3-us-east-1",
			"redshift-us-east-1"
		]
	}

## How does it work
asoi reads the status RSS from AWS and keeps track of the latest update on a
given topic by storing its timestamp in a file named after the topic in the
current directory.

This bot doesn't join the channel, target channels needs to have `-n` mode.
