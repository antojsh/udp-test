var PORT = 33333;
var HOST = '127.0.0.1';

var dgram = require('dgram');


var client = dgram.createSocket('udp4');
var i=0;
function makeid()
{
    var text = "";
    var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

    for( var i=0; i < 5; i++ )
        text += possible.charAt(Math.floor(Math.random() * possible.length));

    return text;
}
setInterval(function(){ 
	i++;
	var message = new Buffer( makeid()+'REV001911176215+1104711-0726929402529532;');
	client.send(message, 0, message.length, PORT, HOST, function(err, bytes) {
	    if (err) throw err;
	    console.log('UDP message sent to ' + HOST +':'+ PORT);
	});
}, 1 );
	

 //client.close();