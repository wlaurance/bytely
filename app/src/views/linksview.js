"use strict";

var $ = require("jquery");
var Backbone = require("backbone");
Backbone.$ = $;
var _ = require("underscore");

var App = require("../app");
var Links = require("../collections/links");
var Link = require("../models/link");
var Templates = require("../templates");

// The LinksView controls the display of a user's LinksCollection. Each
// link can be expanded into a "more info" state. The view deals with this
// by appending expanded links to an array and removing them when they 
// are collapsed. 
module.exports = Backbone.View.extend({

    templateCollapsed: _.template(Templates.urlcollapsed), 
    templateExpanded: _.template(Templates.urlexpanded),

    // Contains ids of expanded links. 
    expandedLinks: [],

    events: {
        "click .urlcollapsed": "toggleCollapsed"
    },

    initialize: function() {
        this.links = new Links();
        this.links.fetch({ silent: true, async: false });
        this.render();
    },

    render: function() {
        var self = this;

        // Add each link to the DOM. 
        this.links.each(function(link) {
            self.$el.prepend(self.templateCollapsed({ link: link }));
        });

        // Update every time a new link is saved. 
        this.listenTo(this.links, "add", this.update);
    },

    update: function() {
        this.$el.empty();
        this.unbind();
        this.render();
    },

    add: function(attributes) {
        var link = new Link(attributes);
        link.save();
        this.links.add(link);
    },

    toggleCollapsed: function(e) {
        var target = $(e.currentTarget);
        if (_.contains(this.expandedLinks, target.data("hash"))) 
            this.collapseLink(target);
        else
            this.expandLink(target);
    },

    collapseLink: function(element) {
        // Remove the additional info box from the DOM. 
        element.removeClass("selected");
        element.next().remove();

        // Remove id from state array. 
        this.expandedLinks = _.reject(this.expandedLinks, function(hash) {
            return hash == element.data("hash");
        });
    },

    expandLink: function(element) {
        // Find the link model that we need. 
        var targetLink = this.links.find(function(link) {
            return link.get("hash") === element.data("hash");
        });        

        element.addClass("selected");
        // Insert the expanded link template into the DOM. 
        element.after(this.templateExpanded({ link: targetLink, appUrl : App.appURL }));

        // Save the state of the link. 
        this.expandedLinks.push(element.data('hash'));
    },

});