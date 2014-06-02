'use strict';

var $ = require('jquery');
var Backbone = require('backbone');
Backbone.$ = $;
var _ = require('underscore');

var Router = require('./router');
var Session = require('./models/session');

// Load form validation helpers.
require("./helpers");

// Make require()ing this file export our app object.
var App = module.exports;
App.root = "/";
App.apiURL = "http://localhost:1337";
App.appURL = "http://localhost:7331";

// Retrieve the user's cookie if they have one. 
App.Session = new Session();
App.Session.initialize();

var backboneSync = Backbone.sync;

// Override Backbone.sync to reference our API endpoint. 
Backbone.sync = function(method, model, options) {
    // Set the API route. 
    options = _.extend(options, {
        url: App.apiURL + (_.isFunction(model.url) ? model.url() : model.url)
    });

    // Add the user's API token to request header. 
    if (App.Session.authenticated()) {
        options.headers = { 'X-Access-Token': App.Session.get('token') }
    }

    backboneSync(method, model, options);
};

// Kick off the app. 
App.router = new Router();
Backbone.history.start({ root: App.root });