"use strict";

var Backbone = require("backbone");
var $ = require('jquery');
var _ = require("underscore");
var Hokey = require("hokeypokey");

var App = require("../app");
var Templates = require("../templates");

Backbone.$ = $;

module.exports = Backbone.View.extend({

    el: 'nav',

    template: _.template(Templates.nav),
    registerTemplate: _.template(Templates.registerForm), 
    signInTemplate: _.template(Templates.signInForm),

    // State variables that determine which form to show.
    isRegistering: false, 
    isSigningIn: false, 

    events: {
        "click #registerlink": "toggleRegisterForm", 
        "click #signinlink": "toggleSignInForm", 
        "click #logoutlink": "logout",
        "click #submitregister": "submitRegisterForm",
        "click #submitsignin": "submitSignInForm"
    },

    initialize: function() {
        this.render();
    },

    render: function() {
        this.$el.html(this.template({ Session: App.Session }));

        if (this.isRegistering) 
            this.$el.append(this.registerTemplate());
        else if (this.isSigningIn) 
            this.$el.append(this.signInTemplate());
    },

    update: function() {
        this.$el.empty();
        this.unbind();
        this.render();
    },

    toggleRegisterForm: function() {
        this.isSigningIn = false;
        this.isRegistering = !this.isRegistering;
        this.update();
    },

    toggleSignInForm: function() {
        this.isRegistering = false;
        this.isSigningIn = !this.isSigningIn;
        this.update();
    },

    logout: function() {
        App.Session.clear();
    },

    submitRegisterForm: function(e) {
        e.preventDefault();

        var form = Hokey(this.$el).rules(function() {
            this.rule("email", /.+\@.+\..+/, "invalidTextField");
            this.rule("password", /\w+/, "invalidTextField");
            this.rule("passwordConfirm", /\w+/, "invalidTextField");
        });

        if (form.get("password") != form.get("passwordconfirm")) return;

        if (form.validate()) {
            var self = this;
            $.post(App.apiURL + "/users", JSON.stringify({
                email: form.get("email"), 
                password: form.get("password")
            }),
            function(data) {
                App.Session.save(data.token);
            });
        }
    },

    submitSignInForm: function(e) {
        e.preventDefault();

        var form = Hokey(this.$el).rules(function() {
            this.rule("email", /.+\@.+\..+/, "invalidTextField");
            this.rule("password", /\w+/, "invalidTextField");
        });

        if (form.validate()) {
            var self = this;
            $.get(App.apiURL + "/token", {
                email: form.get("email"),
                password: form.get("password")
            },
            function(data) {
                App.Session.save(data.token);
            });
        }
    }

});