"use strict";

var $ = require("jquery");
var Backbone = require("backbone"); 
Backbone.$ = $;

var App = require("./app");
var AppView = require("./views/appview");

module.exports = Backbone.Router.extend({

    routes: {
        "": "index"
    },

    index: function() {
        new AppView().render();
    }

})