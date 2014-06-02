var Hokey = require("hokeypokey");

Hokey.helpers(function() {

    this.register("invalidTextField", function(valid, selector) {
        selector.removeClass("error");
        if (!valid) {
            selector.addClass("error");
        }
    });

});