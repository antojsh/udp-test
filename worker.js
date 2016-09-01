self.addEventListener('message',function(e){
	 socket.on('news', function (data) {
		self.postMessage(data)
	});
},false)