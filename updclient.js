var PORT = 33333;
var HOST = '127.0.0.1';

var dgram = require('dgram');


var client = dgram.createSocket('udp4');
var i=0;
setInterval(function(){ 
	i++;
	var message = new Buffer('My KungFu is Good! '+i);
	client.send(message, 0, message.length, PORT, HOST, function(err, bytes) {
	    if (err) throw err;
	    console.log('UDP message sent to ' + HOST +':'+ PORT);
	});
}, 10);
	

 //client.close();