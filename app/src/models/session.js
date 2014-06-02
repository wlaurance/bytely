var $ = require("jquery");
var Backbone = require("backbone");
Backbone.$ = $;

// The Session model is a simple (and insecure) method of creating persistent
// user sessions. The API token is stored unencrypted (;_;) in a cookie and the 
// existence of a session is derived from that. 
module.exports = Backbone.Model.extend({
    defaults: { token: null },

    initialize: function() {
        this.set("token", getCookie("token"));
    },

    authenticated: function() {
        return this.get("token") !== null;
    },

    save: function(token) {
        this.set("token", token);
        setCookie("token", this.get("token"), 30);
    },

    clear: function() {
        this.set("token", null);
        setCookie("token", '', -1);
    }
});

function setCookie(name, value, days) {
    var d = new Date();
    d.setTime(d.getTime() + days * 24 * 60 * 60 * 1000);
    var expires = "expires=" + d.toGMTString();
    var cookie = name + '=' + value + "; " + expires;
    document.cookie = cookie;
}

function getCookie(name) {
    var name = name + '=';
    var cookies = document.cookie.split(';');
    for (var i = 0; i < cookies.length; i++) {
        var cookie = cookies[i].trim();
        if (cookie.indexOf(name) == 0) {
            return cookie.substring(name.length, cookie.length);
        }
    }
    return null;
}