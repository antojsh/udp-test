'use strict'
const PORT = 33333;
const HOST = '127.0.0.1';
const mongoose = require('mongoose');
const dgram = require('dgram');
const server = dgram.createSocket('udp4');
const datagramsModel= require('./models/datagrams')

const client = dgram.createSocket('udp4');
const async = require('async')
const app = require('http').createServer()
const io = require('socket.io')(app);
const fs = require('fs');
const socketConnect = [];
const telnet = require('telnet')
const net = require('net')
var telnetConects=[]


mongoose.Promise = global.Promise;
mongoose.connect('mongodb://localhost/udpdatagram');

server.on('listening', function () {
	let address = server.address();
	console.log('UDP Server listening on ' + address.address + ":" + address.port);
});

server.on('message', function (message, remote) {
	//console.log(String(message))
	async.series([
		saveMongo,
		sendWeb,
		repliTelnet
	])
	function saveMongo(done){
		let datagram = new datagramsModel({
			datagram:String(message),
			timestamp:Date.now()
		})
		datagram.save(function(err,datagramresponse){
			if(err)
				console.log('No guardado')	
		})
		done()
	}

	function sendWeb(done){
		io.sockets.emit('news', String(message));
		done()
	}

	function repliTelnet(){
			for (var i = telnetConects.length - 1; i >= 0; i--) {
				telnetConects[i].write(String(message))
			}
	}
});


io.on('connection', function (socket) {
	socket.on('enviar evento', function (data) {
	    console.log(data);
	});
});
net.createServer(function(sock) {
    telnetConects.push(sock)
    // We have a connection - a socket object is assigned to the connection automatically
    console.log('CONNECTED: ' + sock.remoteAddress +':'+ sock.remotePort);
    
    // Add a 'data' event handler to this instance of socket
    sock.on('data', function(data) {
        
        console.log('DATA ' + sock.remoteAddress + ': ' + data);
        // Write the data back to the socket, the client will receive it as data from the server
        sock.write('You said "' + data + '"');
        
    });
    
    // Add a 'close' event handler to this instance of socket
    sock.on('close', function(data) {
    	telnetConects.splice(telnetConects.indexOf(sock), 1);
        console.log('CLOSED: ' + sock.remoteAddress +' '+ sock.remotePort);
    });
    
}).listen('8000','192.168.129.178');


server.bind(PORT, HOST);//UDP SERVER
app.listen(3000);// SOCKET IO SERVER


