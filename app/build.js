var fs = require('fs');
var sys = require('sys');
var exec = require('child_process').exec;

var _ = require('underscore');

// Create a giant object full of all of the templates in the
// src/templates directory.
var directory = 'src/templates/';
var compiled = {};
var filenames = fs.readdirSync(directory);

_.each(filenames, function(file) {
    var templateKey = file.split('.')[0];
    var template = fs.readFileSync(directory + file, {
        encoding: 'utf8'
    });
    compiled[templateKey] = template;
});

// Write this object out to a file that can be included via
// Browserify. 
var output = "module.exports = " + JSON.stringify(compiled) + ';';
fs.writeFileSync('src/templates.js', output);

// Build the bundle.
exec("browserify src/app.js -o bundle.js", function(err, stdout, stderr) {
    sys.puts(stdout);
    sys.puts(stderr);
});