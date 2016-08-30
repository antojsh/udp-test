'use strict'
var PORT = 33333;
var HOST = '127.0.0.1';
var mongoose = require('mongoose');
var dgram = require('dgram');
var server = dgram.createSocket('udp4');
var datagramsModel= require('./models/datagrams')
var dgram = require('dgram');
var client = dgram.createSocket('udp4');

var app = require('http').createServer()
var io = require('socket.io')(app);
var fs = require('fs');
var socketConnect = [];
mongoose.Promise = global.Promise;
mongoose.connect('mongodb://localhost/udpdatagram');

server.on('listening', function () {
	var address = server.address();
	console.log('UDP Server listening on ' + address.address + ":" + address.port);
});

server.on('message', function (message, remote) {
	console.log(message)
	let datagram = new datagramsModel({
		datagram:message,
		timestamp:Date.now()
	})
	datagram.save(function(err,datagramresponse){
		if(err)
			console.log('No guardado')	
	})
	
	sendWeb(message)
});

function sendWeb(message){
	io.sockets.emit('news', String(message));
}

io.on('connection', function (socket) {});
server.bind(PORT, HOST);
app.listen(3000);


