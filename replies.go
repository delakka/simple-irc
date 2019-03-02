package main

type Reply struct {
	Numeric string
	Message string
}

var (
	RPL_WELCOME = Reply{
		Numeric: "001",
		Message: ":Welcome to the server!",
	}
	RPL_TOPIC = Reply{
		Numeric: "332",
	}
	RPL_NAMREPLY = Reply{
		Numeric: "353",
	}
	RPL_ENDOFNAMES = Reply{
		Numeric: "366",
		Message: ":End of NAMES list",
	}
	ERR_NORECIPIENT = Reply{
		Numeric: "411",
		Message: ":No recipients were found",
	}
	ERR_NOTEXTTOSEND = Reply{
		Numeric: "412",
		Message: ":No message was given",
	}
	ERR_NONICKNAMEGIVEN = Reply{
		Numeric: "431",
		Message: ":No nick was given",
	}
	ERR_NICKNAMEINUSE = Reply{
		Numeric: "433",
		Message: ":Nick is already in use",
	}
	ERR_NICKCOLLISION = Reply{
		Numeric: "436",
	}
	ERR_NOTONCHANNEL = Reply{
		Numeric: "442",
		Message: ":The user is not in the specified channel",
	}
	ERR_NEEDMOREPARAMS = Reply{
		Numeric: "461",
		Message: ":Need more parameters",
	}
	//RPL_AWAY             = "301"
	//ERR_NOSUCHNICK       = "401"
	//ERR_NOSUCHSERVER     = "402"
	//ERR_NOSUCHCHANNEL    = "403"
	//ERR_CANNOTSENDTOCHAN = "404"
	//ERR_TOOMANYCHANNELS  = "405"
	//ERR_TOOMANYTARGETS   = "407"
	//ERR_NOTOPLEVEL       = "413"
	//ERR_WILDTOPLEVEL     = "414"
	//ERR_TOOMANYMATCHES   = "416"
	//ERR_ERRONEUSNICKNAME = "432"
	//ERR_UNAVAILRESOURCE  = "437"
	//ERR_ALREADYREGISTRED = "462"
	//ERR_CHANNELISFULL    = "471"
	//ERR_INVITEONLYCHAN   = "473"
	//ERR_BANNEDFROMCHAN   = "474"
	//ERR_BADCHANNELKEY    = "475"
	//ERR_BADCHANMASK      = "476"
	//ERR_RESTRICTED       = "484"
)
