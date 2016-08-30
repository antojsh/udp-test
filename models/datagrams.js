const mongoose = require('mongoose');
const Schema 	 = mongoose.Schema;
const datagrams_schema = new Schema({
	datagram: {
		type:String
	},
	timestamp: {
		type:Date
	}
});

var datagrams = mongoose.model('datagrams',datagrams_schema);
module.exports= datagrams;