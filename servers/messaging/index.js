//put interpreter into strict mode
// "use strict";

//require the express and morgan packages
const express = require("express");
const morgan = require("morgan");

//create a new express application
const app = express();

//get ADDR environment variable,
//defaulting to ":80"
// const addr = "localhost:4000";
const addr = process.env.ADDR || ":80";
//split host and port using destructuring
const [host, port] = addr.split(":");

const channel = require('./modelChannel');
const message = require('./modelMessage');

var mysql = require('mysql');
var amqp = require('amqplib/callback_api');

var pool = mysql.createPool({ // same values as go
    host     : process.env.MYSQL_ADDR,
    user     : "root",
    password : process.env.MYSQL_ROOT_PASSWORD,
    database : process.env.MYSQL_DB
  });

var rabbitChannel;
amqp.connect('amqp://' + process.env.RABBIT + ":5672/", function(err, conn) {
    conn.createChannel(function(err, ch) {
        var q = process.env.RABBIT;

        ch.assertQueue(q, {durable: false});
        // Note: on Node 6 Buffer.from(msg) should be used
        rabbitChannel = ch;
    });
    // setTimeout(function() { conn.close(); process.exit(0) }, 500);
});

//add JSON request body parsing middleware
app.use(express.json());
//add the request logging middleware
app.use(morgan("dev"));

app.get("/v1/channels", (req, res, next) => {
    let user = req.get("X-User")
    if (!user) {
        res.status(401).send({error: 'ChannelsHandler: User unauthorized'});
        return
    }
    user = JSON.parse(user)

    let channels = [];
    pool.query('select * from channels', function (error, results, fields) {
        if (error) { 
            res.status(404).send("No Channels Found/Error occured when finding channels")
            return
        }
        if (results.length != 0) {
            results.forEach(element => {
                channels.push(element);
            });
        }
        
        res.setHeader("Content-Type", "application/json");
        try {
            res.status(200).json(channels);
            return channels;
        } catch(err) {
            res.status(500).send("Error encoding channels into json")
            next(err);
        }
    });
});
    
app.post("/v1/channels", (req, res, next) => {
    let user = req.get("X-User")
    if (!user) {
        res.status(401).send({error: 'ChannelsHandler: User unauthorized'});
        return
    }
    user = JSON.parse(user)
    // console.log("Req body" + req.body.nameString + " " + req.body.descriptionString + " " + req.body.privateBool + " " + user.id)
    let newChannel = new channel(req.body.nameString, req.body.descriptionString, req.body.privateBool, user.id)
   
    pool.query('insert into channels (nameString, descriptionString, privateBool, createdAt, creatorID, editedAt) values (?, ?, ?, ?, ?, ?)', 
    [newChannel.nameString, newChannel.descriptionString, newChannel.privateBool, newChannel.createdAt, newChannel.creatorID, newChannel.editedAt], 
    function (error, results, fields) {
        if (error || results.affectedRows == 0) { 
            console.log(error)
            res.status(500).send("Can't Insert Channel/Error occured when inserting channel")
            return
        }
        newChannel.id = results.insertId
        pool.query('insert into channels_members (channelID, userID) values (?,?)', 
            [newChannel.id, newChannel.creatorID],
            function (error, results, fields) {
                // console.log(newChannel.id, newChannel.creatorID)
                if (error || results.affectedRows == 0) { 
                    res.status(500).send("Can't Insert Member/Error occured when inserting member into channel_member")
                    return
                }   
             
                res.setHeader("Content-Type", "application/json")
                let channelMembersEncoded;
                if (channel.privateBool) {
                    pool.query('select * from channels_members where channelID = ?', channel.id, function (error, results, fields) {
                        if (error) { 
                            res.status(404).send("No Channels Found/Error occured when finding channels")
                            return
                        }
                        if (results.length != 0) {
                            results.forEach(element => {
                                channels.push(element);
                            });
                        }
                        channelMembersEncoded = results;
                        // channelMembersEncoded = JSON.stringify(results)
                    });
                }

                try {
                    // j, err = json.Marshal(newChannel)
                    // if (err != nil) {
                    //     fmt.Sprintf("Error encoding new channel %v", err)
                    //     return
                    // }
                    var newEvent = {
                        // always initialize all instance properties
                        type: "channel-new",
                        channel: newChannel,
                        userIDs: channelMembersEncoded
                    }
                    
                    // Note: on Node 6 Buffer.from(msg) should be used
                    rabbitChannel.sendToQueue(process.env.RABBIT, new Buffer(JSON.stringify(newEvent)));
                    console.log(" [x] Sent 'New Channel!'");
                    
                    res.status(201).json(newChannel);
                    res.end();
                    return newChannel;
                } catch(err) {
                    res.status(500).send("Error encoding")
                    next(err);
                }
        });
    });
});

app.get("/v1/channels/:chanid", (req, res, next) => {
    let user = req.get("X-User")
    if (!user) {
        res.status(401).send({error: 'ChannelsHandler: User unauthorized'});
        return
    }
    user = JSON.parse(user)

    let channelID = req.params.chanid;
    pool.query('select * from channels where id = ?', channelID, function (error, results, fields) {
        if (error) { 
            res.status(500).send("Can't Get Channel by ID")
            return
        }
        if (results.length == 0) {
            res.status(200).json(results);
            return
        }

        let channel = results[0]
        let messages100 = []
        let messageBefore = []

        if (channel.nameString == "general") {
            if (Object.keys(req.query).length != 0) {
                pool.query('select * from messages where id < ? and channelID = ? order by createdAt desc limit 100', [req.query.before, channel.id], function (error, results, fields) {
                    if (error) { 
                        res.status(500).send("Can't get before top 100 messages")
                        return 
                    }
                    
                    if (results.length != 0) {
                        results.forEach(element => {
                            messageBefore.push(element);
                        });
                    } 
                    res.setHeader("Content-Type", "application/json");
                    try {
                        res.status(200).json(messageBefore);
                        res.end();
                        return messageBefore;
                    } catch(err) {
                        res.status(500).send("Error encoding");
                        next(err);
                    }
                });
            } else {
                pool.query('select * from messages where channelID = ? limit 100', channelID, 
                function (error, results, fields) {
                    if (error) { 
                        res.status(500).send("Can't get Top 100 Messages from channel")
                        return 
                    }
                    
                    if (results.length != 0) {
                        results.forEach(element => {
                            messages100.push(element);
                        });
                    }

                    res.setHeader("Content-Type", "application/json");
                    try {
                        res.status(200).json(messages100);
                        res.end();
                        return messages100;
                    } catch(err) {
                        res.status(500).send("Error encoding");
                        next(err);
                    }
                });
            }
        } else {
            pool.query('select * from channels_members where userID = ? and channelID = ?', [user.id, channelID], function (error, results, fields) {
                if (error) {
                    res.status(500).send("Can't get channel members when trying to authorize")
                    return 
                }

                if (channel.privateBool && results.length == 0) { 
                    res.status(403).send("current user was not invited into the private channel");
                    return;
                }
                
                if (Object.keys(req.query).length != 0) {
                    pool.query('select * from messages where id < ? and channelID = ? order by createdAt desc limit 100', [req.query.before, channel.id], function (error, results, fields) {
                        if (error) { 
                            res.status(500).send("Can't get before top 100 messages")
                            return 
                        }
                        
                        if (results.length != 0) {
                            results.forEach(element => {
                                messageBefore.push(element);
                            });
                        } 
                        res.setHeader("Content-Type", "application/json");
                        try {
                            res.status(200).json(messageBefore);
                            res.end();
                            return messageBefore;
                        } catch(err) {
                            res.status(500).send("Error encoding");
                            next(err);
                        }
                    });
                } else {
                    pool.query('select * from messages where channelID = ? limit 100', channelID, 
                    function (error, results, fields) {
                        if (error) { 
                            res.status(500).send("Can't get Top 100 Messages from channel")
                            return 
                        }
                        
                        if (results.length != 0) {
                            results.forEach(element => {
                                messages100.push(element);
                            });
                        }

                        res.setHeader("Content-Type", "application/json");
                        try {
                            res.status(200).json(messages100);
                            res.end();
                            return messages100;
                        } catch(err) {
                            res.status(500).send("Error encoding");
                            next(err);
                        }
                    });
                }
            }); 
        } 
    });
});
  
app.post("/v1/channels/:chanid", (req, res, next) => {
    let user = req.get("X-User")
    if (!user) {
        res.status(401).send({error: 'ChannelsHandler: User unauthorized'});
        return
    }
    user = JSON.parse(user)
    
    let channelID = req.params.chanid;

    pool.query('select * from channels where id = ?', channelID, function (error, results, fields) {
        if (error) { 
            res.status(500).send("Can't Get Channel by ID")
            return
        }
        if (results.length == 0) {
            res.status(500).send("Specified Channel doesn't exist")
            return
        }
        let channel = results[0]
        // console.log(channel)
        let channelMembersEncoded;
        if (channel.id == 1) {
            let newMessage = new message(channel.id, req.body.body, user.id)
           
            pool.query('insert into messages (channelID, body, createdAt, creatorID, editedAt) values (?, ?, ?, ?, ?)', 
            [newMessage.channelID, newMessage.body, newMessage.createdAt, newMessage.creatorID, newMessage.editedAt], 
            function (error, results, fields) {
                // console.log(results)
                newMessage.id = results.insertId
                if (error || results.affectedRows == 0) { 
                    // console.log(error)
                    res.status(500).send("Can't insert messages into table")
                    return
                }
    
                res.setHeader("Content-Type", "application/json");
                
                try {
                    if (channel.privateBool) {
                        pool.query('select * from channels_members where channelID = ?', channel.id, function (error, results, fields) {
                            if (error) { 
                                res.status(404).send("No Channels Found/Error occured when finding channels")
                                return
                            }
                            // if (results.length != 0) {
                            //     results.forEach(element => {
                            //         channels.push(element);
                            //     });
                            // }
                            channelMembersEncoded = results;
                        });
                    }
    
                    // j, err = json.Marshal(newMessage)
                    // if (err != nil) {
                    //     fmt.Sprintf("Error encoding new message %v", err)
                    //     return
                    // }
                    var newEvent = {
                        // always initialize all instance properties
                        type: "message-new",
                        message: newMessage,
                        userIDs: channelMembersEncoded
                    }
                    // Note: on Node 6 new Buffer(JSON.stringify(newEvent))(msg) should be used
                    rabbitChannel.sendToQueue(process.env.RABBIT, new Buffer(JSON.stringify(newEvent)));
                    console.log(" [x] Sent 'New Message!'");
                    res.status(201).json(newMessage);
                    res.end();
                    return newMessage;
                } catch(err) {
                    res.status(500).send("Error encoding messages into json");
                    next(err);
                }
            });
        } else {
            pool.query('select * from channels_members where userID = ? and channelID = ?', [user.id, channelID], function (error, results, fields) {
                if (error) {
                    res.status(500).send("Can't get channel members when trying to authorize")
                    return 
                }
                
                if (channel.privateBool && results.length == 0) { 
                    res.status(403).send("current user was not invited into the private channel");
                    return;
                }
                console.log(results)
                // let newMessage = new message(channel.id, req.body.body, user.id)
                let newMessage = new message(results[0].channelID, req.body.body, results[0].userID)
            
                pool.query('insert into messages (channelID, body, createdAt, creatorID, editedAt) values (?, ?, ?, ?, ?)', 
                [newMessage.channelID, newMessage.body, newMessage.createdAt, newMessage.creatorID, newMessage.editedAt], 
                function (error, results, fields) {
                    // console.log(results)
                    newMessage.id = results.insertId
                    if (error || results.affectedRows == 0) { 
                        console.log(error)
                        res.status(500).send("Can't insert messages into table")
                        return
                    }
        
                    res.setHeader("Content-Type", "application/json");
                   
                    try {
                        if (channel.privateBool) {
                            pool.query('select * from channels_members where channelID = ?', channel.id, function (error, results, fields) {
                                if (error) { 
                                    res.status(404).send("No Channels Found/Error occured when finding channels")
                                    return
                                }
                                // if (results.length != 0) {
                                //     results.forEach(element => {
                                //         channels.push(element);
                                //     });
                                // }
                                channelMembersEncoded = results
                            });
                        }
        
                        // j, err = json.Marshal(newMessage)
                        // if (err != nil) {
                        //     fmt.Sprintf("Error encoding new message %v", err)
                        //     return
                        // }
                        var newEvent = {
                            // always initialize all instance properties
                            type: "message-new",
                            message: newMessage,
                            userIDs: channelMembersEncoded
                        }
                        // Note: on Node 6 Buffer.from(msg) should be used
                        rabbitChannel.sendToQueue(process.env.RABBIT, new Buffer(JSON.stringify(newEvent)));
                        console.log(" [x] Sent 'New Message!'");
                        res.status(201).json(newMessage);
                        res.end();
                        return newMessage;
                    } catch(err) {
                        res.status(500).send("Error encoding messages into json");
                        next(err);
                    }
                });
            }); 
        }
    });
});

app.patch("/v1/channels/:chanid", (req, res, next) => {
    let user = req.get("X-User")
    if (!user) {
        res.status(401).send({error: 'ChannelsHandler: User unauthorized'});
        return
    }
    user = JSON.parse(user)
    
    let channelID = req.params.chanid;

    pool.query('select * from channels where id = ?', channelID, function (error, results, fields) {
        if (error) { 
            res.status(500).send("Can't Get Channel by ID")
            return
        }
        if (results.length == 0) {
            res.status(500).send("Specified Channel doesn't exist for patch")
            return
        }
        let channel = results[0]
        if (channel.nameString != "general") {
            pool.query('select * from channels_members where userID = ? and channelID = ?', [user.id, channelID], function (error, results, fields) {
                if (error) {
                    res.status(500).send("Can't get channel members when trying to authorize")
                    return 
                }
                
                if (channel.privateBool && results.length == 0) { 
                    res.status(403).send("current user was not invited into the private channel");
                    return;
                }

                pool.query('update channels set nameString = ?, descriptionString = ? where id = ?', [req.body.nameString, req.body.descriptionString, channelID], 
                function (error, results, fields) {
                    if (error || results.affectedRows == 0) {
                        res.status(500).send("Can't update channels")
                        return
                    }
                    pool.query('select * from channels where id = ?', channelID, function (error, results, fields) {
                        if (error) { 
                            res.status(500).send("Can't Get Channel by ID")
                            return
                        }
                        if (results.length == 0) {
                            res.status(500).send("Specified Channel doesn't exist for patch")
                            return
                        }
                        
                        res.header("Content-Type", "application/json");

                        let channelMembersEncoded;
                        if (channel.privateBool) {
                            pool.query('select * from channels_members where channelID = ?', channel.id, function (error, results, fields) {
                                if (error) { 
                                    res.status(404).send("No Channels Found/Error occured when finding channels")
                                    return
                                }
                                // if (results.length != 0) {
                                //     results.forEach(element => {
                                //         channels.push(element);
                                //     });
                                // }
                                channelMembersEncoded = results;
                            });
                        }

                        try {
                            // j, err = json.Marshal(newChannel)
                            // if (err != nil) {
                            //     fmt.Sprintf("Error encoding new channel %v", err)
                            //     return
                            // }
                            var newEvent = {
                                // always initialize all instance properties
                                type: "channel-update",
                                channel: results[0],
                                userIDs: channelMembersEncoded
                            }
                            // Note: on Node 6 Buffer.from(msg) should be used
                            rabbitChannel.sendToQueue(process.env.RABBIT, new Buffer(JSON.stringify(newEvent)));
                            console.log(" [x] Updated New Channel");
                            res.status(200).json(results[0]);
                            res.end();
                            return results[0];
                        } catch(err) {
                            res.status(500).send("Error encoding new channel into json");
                            next(err);
                        }
                    });
                });
            }); 
        } else {
            res.status(500).send("You are not authorized to change general's name string and/or description")
            return 
        }
    });
});

app.delete("/v1/channels/:chanid", (req, res, next) => {
    let user = req.get("X-User")
    if (!user) {
        res.status(401).send({error: 'ChannelsHandler: User unauthorized'});
        return
    }
    user = JSON.parse(user)

    let channelID = req.params.chanid;
    if (channelID != 1) {
        pool.query('select * from channels where id = ?', channelID, function (error, results, fields) {
            if (error) { 
                res.status(500).send("Can't Get Channel by ID")
                return
            }
            if (results.length == 0) {
                res.status(500).send("Specified Channel doesn't exist for delete")
                return
            }
            let channel = results[0]
            pool.query('select * from channels_members where userID = ? and channelID = ?', [user.id, channelID], function (error, results, fields) {
                if (error) {
                    res.status(500).send("Can't get channel members when trying to authorize")
                    return 
                }
                
                if (channel.privateBool && results.length == 0) { 
                    res.status(403).send("current user was not invited into the private channel");
                    return;
                }

                pool.query('delete from channels where id = ?', channelID, function (error, results, fields) {
                    if (error || results.affectedRows == 0) {
                        res.status(500).send("Can't delete channels")
                        return
                    }

                    pool.query('delete from messages where channelID = ?', channelID, function (error, results, fields) {
                        if (error) {
                            res.status(500).send("Can't delete messages")
                            return
                        }

                        pool.query('delete from channels_members where channelID = ?', channelID, function (error, results, fields) {
                            if (error || results.affectedRows == 0) {
                                res.status(500).send("Can't delete channels_members")
                                return
                            }
               
                            let channelMembersEncoded;
                            if (channel.privateBool) {
                                pool.query('select * from channels_members where channelID = ?', channel.id, function (error, results, fields) {
                                    if (error) { 
                                        res.status(404).send("No Channels Found/Error occured when finding channels")
                                        return
                                    }
                                    // if (results.length != 0) {
                                    //     results.forEach(element => {
                                    //         channels.push(element);
                                    //     });
                                    // }
                                    channelMembersEncoded = results;
                                });
                            }
            
                            var newEvent = {
                                // always initialize all instance properties
                                type: "channel-delete",
                                channelID: channelID,
                                userIDs: channelMembersEncoded
                            }
                            // Note: on Node 6 Buffer.from(msg) should be used
                            rabbitChannel.sendToQueue(process.env.RABBIT, new Buffer(JSON.stringify(newEvent)));
                            console.log(" [x] Sent 'Channel Deleted!'");

                            res.status(200).send("channel deleted successfully");
                            return
                        }); 
                    }); 
                });
            }); 
        });
    } else {
        res.status(500).send("You are not authorized to delete the general channel")
        return 
    }
});

app.post("/v1/channels/:chanid/members", (req, res, next) => {
    let user = req.get("X-User")
    if (!user) {
        res.status(401).send({error: 'ChannelsHandler: User unauthorized'});
        return
    }
    user = JSON.parse(user)

    let channelID = req.params.chanid;
    pool.query('select * from users where id = ?', req.body.id, function (error, results, fields) {
        if (error) { 
            res.status(500).send("Can't Get User by ID")
            return
        }
        if (results.length == 0) {
            res.status(500).send("Specified User doesn't exist in database")
            return
        }
        pool.query('select * from channels where id = ?', channelID, function (error, results, fields) {
            if (error) { 
                res.status(500).send("Can't Get Channel by ID")
                return
            }
            if (results.length == 0) {
                res.status(500).send("Specified Channel doesn't exist for post")
                return
            }
            let channel = results[0]
        
            if (channel.nameString == "general") {
                pool.query('insert into channels_members (channelID, userID) values (?, ?)', [channelID, req.body.id],
                function (error, results, fields) {
                    console.log(error)
                    if (error || results.affectedRows == 0) {
                        res.status(500).send("Can't add to members")
                        return
                    }

                    res.status(201).send("user was added to the channel");
                    return results;
                });
            } else {
                pool.query('select * from channels_members where userID = ? and channelID = ?', [user.id, channelID], function (error, results, fields) {
                    if (error) {
                        res.status(500).send("Can't get channel members when trying to authorize")
                        return 
                    }
                    
                    if (channel.privateBool && results.length == 0) { 
                        res.status(403).send("current user was not invited into the private channel");
                        return;
                    }

                    pool.query('insert into channels_members (channelID, userID) values (?, ?)', [channelID, req.body.id],
                    function (error, results, fields) {
                        console.log(error)
                        if (error || results.affectedRows == 0) {
                            res.status(500).send("Can't add to members")
                            return
                        }

                        res.status(201).send("user was added to the channel");
                        return results;
                    });
                }); 
            }
        });
    });
});

app.delete("/v1/channels/:chanid/members", (req, res, next) => {
    let user = req.get("X-User")
    if (!user) {
        res.status(401).send({error: 'ChannelsHandler: User unauthorized'});
        return
    }
    user = JSON.parse(user)

    let channelID = req.params.chanid;
    if (channelID != 1) {
        pool.query('select * from users where id = ?', req.body.id, function (error, results, fields) {
            if (error) { 
                res.status(500).send("Can't Get User by ID")
                return
            }
            if (results.length == 0) {
                res.status(500).send("Specified User doesn't exist in database")
                return
            }
        
            pool.query('select * from channels where id = ?', channelID, function (error, results, fields) {
                if (error || results.length == 0) { 
                    res.status(500).send("Can't Get Channel by ID")
                    return
                }
                let channel = results[0]
                
                pool.query('select * from channels_members where userID = ? and channelID = ?', [user.id, channelID], function (error, results, fields) {
                    if (error) {
                        res.status(500).send("Can't get channel members when trying to authorize")
                        return 
                    }
                    
                    if (channel.privateBool && results.length == 0) { 
                        res.status(403).send("current user was not invited into the private channel");
                        return;
                    }

                    pool.query('delete from channels_members where id = ?', req.body.id, function (error, results, fields) {
                        if (error || results.affectedRows == 0) {
                            res.status(500).send("Can't delete from members/Given Member did not exist in Channel's Members list")
                            return
                        }

                        res.status(200).send("user was deleted from channel's members list");
                        return results;
                    });
                }); 
            });
        });
    } else {
        res.status(500).send("You are not authorized to delete the general channel's members. Since it is a public channel, everybody has access to it.")
        return 
    }
});


app.patch("/v1/messages/:messageid", (req, res, next) => {
    let user = req.get("X-User")
    if (!user) {
        res.status(401).send({error: 'ChannelsHandler: User unauthorized'});
        return
    }
    user = JSON.parse(user)
    
    let messageID = req.params.messageid;
    let channelMembersEncoded;
    pool.query('select * from messages where id = ?', messageID, function (error, results, fields) {
        if (error || results.length == 0) { 
            res.status(500).send("Can't Get Message by ID")
            return
        }
        if (results[0].channelID == 1) {
            pool.query('update messages set body = ? where id = ?', [req.body.body, messageID], function (error, results, fields) {
                if (error || results.affectedRows == 0) {
                    res.status(500).send("Can't add to members")
                    return
                }
                pool.query('select * from messages where id = ?', messageID, function (error, results, fields) {
                    if (error || results.length == 0) { 
                        res.status(500).send("Can't Get Message by ID")
                        return
                    }
                    res.header("Content-Type", "application/json");
            
                    try {
                        if (channel.privateBool) {
                            pool.query('select * from channels_members where channelID = ?', channel.id, function (error, results, fields) {
                                if (error) { 
                                    res.status(404).send("No Channels Found/Error occured when finding channels")
                                    return
                                }
                                // if (results.length != 0) {
                                //     results.forEach(element => {
                                //         channels.push(element);
                                //     });
                                // }
                                channelMembersEncoded = results;
                            });
                        }
        
                        // j, err = json.Marshal(newMessage)
                        // if (err != nil) {
                        //     fmt.Sprintf("Error encoding new message %v", err)
                        //     return
                        // }
                        var newEvent = {
                            // always initialize all instance properties
                            type: "message-update",
                            message: results,
                            userIDs: channelMembersEncoded
                        }
                        // Note: on Node 6 Buffer.from(msg) should be used
                        rabbitChannel.sendToQueue(process.env.RABBIT, new Buffer(JSON.stringify(newEvent)));
                        console.log(" [x] Update 'Message!'");
                        res.status(200).json(results);
                        res.end();
                        return results
                    } catch(err) {
                        res.status(500).send("Error encoding channel");
                        next(err);
                    }
                });  
            });
        } else {
            pool.query('select * from channels_members where userID = ? and channelID = ?', [user.id, results[0].channelID], function (error, results, fields) {
                if (error) {
                    res.status(500).send("Can't get channel members when trying to authorize")
                    return 
                }
                
                if (results.length == 0) { 
                    res.status(403).send("current user was not invited into the private channel");
                    return;
                }
                
                pool.query('update messages set body = ? where id = ?', [req.body.body, messageID], function (error, results, fields) {
                    if (error || results.affectedRows == 0) {
                        res.status(500).send("Can't add to members")
                        return
                    }
                    pool.query('select * from messages where id = ?', messageID, function (error, results, fields) {
                        if (error || results.length == 0) { 
                            res.status(500).send("Can't Get Message by ID")
                            return
                        }
                        res.header("Content-Type", "application/json");
                
                        try {
                            if (channel.privateBool) {
                                pool.query('select * from channels_members where channelID = ?', channel.id, function (error, results, fields) {
                                    if (error) { 
                                        res.status(404).send("No Channels Found/Error occured when finding channels")
                                        return
                                    }
                                    // if (results.length != 0) {
                                    //     results.forEach(element => {
                                    //         channels.push(element);
                                    //     });
                                    // }
                                    channelMembersEncoded = results;
                                });
                            }
            
                            // j, err = json.Marshal(newMessage)
                            // if (err != nil) {
                            //     fmt.Sprintf("Error encoding new message %v", err)
                            //     return
                            // }
                            var newEvent = {
                                // always initialize all instance properties
                                type: "message-update",
                                message: results,
                                userIDs: channelMembersEncoded
                            }
                            // Note: on Node 6 Buffer.from(msg) should be used
                            rabbitChannel.sendToQueue(process.env.RABBIT, new Buffer(JSON.stringify(newEvent)));
                            console.log(" [x] Update 'New Message!'");
                            res.status(200).json(results);
                            res.end();
                            return results
                        } catch(err) {
                            res.status(500).send("Error encoding channel");
                            next(err);
                        }
                    });  
                });
            });
        } 
    });
});

app.delete("/v1/messages/:messageid", (req, res, next) => {
    let user = req.get("X-User")
    if (!user) {
        res.status(401).send({error: 'ChannelsHandler: User unauthorized'});
        return
    }
    user = JSON.parse(user)

    let messageID = req.params.messageid;

    pool.query('select * from messages where id = ?', messageID, function (error, results, fields) {
        if (error || results.length == 0) { 
            res.status(500).send("Can't Get Message by ID")
            return
        }
        let channelMembersEncoded;
        if (results[0].channelID == 1) {
            pool.query('delete from messages where id = ?', messageID, 
            function (error, results, fields) {
                if (error || results.affectedRows == 0) {
                    res.status(500).send("Can't add to members")
                    return
                }
    
                res.status(200).send("Deleting the message was successful");

                if (channel.privateBool) {
                    pool.query('select * from channels_members where channelID = ?', channel.id, function (error, results, fields) {
                        if (error) { 
                            res.status(404).send("No Channels Found/Error occured when finding channels")
                            return
                        }
                        // if (results.length != 0) {
                        //     results.forEach(element => {
                        //         channels.push(element);
                        //     });
                        // }
                        channelMembersEncoded = results;
                    });
                }

                var newEvent = {
                    // always initialize all instance properties
                    type: "message-delete",
                    messageID: messageID,
                    userIDs: channelMembersEncoded
                }
                // Note: on Node 6 Buffer.from(msg) should be used
                rabbitChannel.sendToQueue(process.env.RABBIT, new Buffer(JSON.stringify(newEvent)));
                console.log(" [x] Deleted a Message!'");
                res.end();
                return results;
            });
        } else {
            pool.query('select * from channels_members where userID = ? and channelID = ?', [user.id, results[0].channelID], function (error, results, fields) {
                if (error) {
                    res.status(500).send("Can't get channel members when trying to authorize")
                    return 
                }
                
                if (results.length == 0) { 
                    res.status(403).send("current user was not invited into the private channel");
                    return;
                }
                
                pool.query('delete from messages where id = ?', messageID, 
                function (error, results, fields) {
                    if (error || results.affectedRows == 0) {
                        res.status(500).send("Can't add to members")
                        return
                    }
        
                    if (channel.privateBool) {
                        pool.query('select * from channels_members where channelID = ?', channel.id, function (error, results, fields) {
                            if (error) { 
                                res.status(404).send("No Channels Found/Error occured when finding channels")
                                return
                            }
                            // if (results.length != 0) {
                            //     results.forEach(element => {
                            //         channels.push(element);
                            //     });
                            // }
                            channelMembersEncoded = results;
                        });
                    }
    
                    var newEvent = {
                        // always initialize all instance properties
                        type: "message-delete",
                        messageID: messageID,
                        userIDs: channelMembersEncoded
                    }
                    // Note: on Node 6 Buffer.from(msg) should be used
                    rabbitChannel.sendToQueue(process.env.RABBIT, new Buffer(JSON.stringify(newEvent)));
                    console.log(" [x] Deleted a 'New Message!'");
                    res.status(200).send("Deleting the message was successful");
                    res.end();
                    return results;
                });
            }); 
        }
    });
});

//start the server listening on host:port
app.listen(port, host, () => {
  //callback is executed once server is listening
  console.log(`server is listening at http://${addr}...`);
});

