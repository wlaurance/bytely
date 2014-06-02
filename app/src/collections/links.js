'use strict';

var Backbone = require('backbone');
var LinkModel = require('../models/link');

module.exports = Backbone.Collection.extend({
    model: LinkModel, 
    url: '/links'
});