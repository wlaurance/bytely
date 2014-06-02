"use strict";

var $ = require("jquery");
var Backbone = require("backbone");
Backbone.$ = $;
var _ = require("underscore");

var App = require("../app");
var Templates = require("../templates");
var NavView = require("./navview");
var LinksView = require("./linksview");

// The AppView is the most abstract view in the app. It directly controls the url
// shortening form and has subviews to control the navigation/login forms and the 
// display/manipulation of the Links collection. 
module.exports = Backbone.View.extend({

    el: "body",

    template: _.template(Templates.app),

    events: {
        "click #submiturl": "submitUrl"
    },

    initialize: function() {},

    render: function() {
        this.$el.prepend(this.template({ Session: App.Session }));

        this.formview = new NavView();
        this.linksview = new LinksView({ el: "#linksview" });

        // Refresh page on change of session state.
        this.listenTo(App.Session, "change", this.update);

        return this;
    },

    update: function() {
        this.$el.empty();
        this.unbind();
        this.render();
    },

    submitUrl: function() {
        var link = this.$el.find("#newurlbox").val();
        if (link) this.linksview.add({ original_url: link });
    }
    
});