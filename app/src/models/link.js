"use strict";

var Backbone = require("backbone");

module.exports = Backbone.Model.extend({

    url: "/links", 
    
    idAttribute: "hash",

    defaults: {
        original_url: "", 
        created_on: currentDate(),
        hits: 0, 
        last_hit: "",
        mobile_hits: 0
    }, 

    initialize: function(attributes, options) {
        
    },

    // As the unique id (the hash) is generated on the server, if the model is
    // new we must intercept the response and manually set the ID.
    save: function(attributes, options) {
        var self = this;
        var opts = options || {}; 
        if (this.isNew()) {
            opts.success = function(model, response) {
                self.id = response.hash;
            }
            opts.async = false;
        }
        Backbone.Model.prototype.save.call(this, attributes, opts);
    }

});

function currentDate() {
    var date = new Date();
    return date.getMonth() + '/' + date.getDate() + '/' + date.getFullYear();
}