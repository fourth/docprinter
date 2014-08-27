var page = require('webpage').create(),
    system = require('system'),
    address, output, size;

if (system.args.length !=2) {
    console.log('Usage: renderpdf.js URL');
    phantom.exit(1);
} else {
    address = system.args[1];
    output = "/dev/stdout";
    page.viewportSize = { width: 600, height: 600 };
    page.paperSize = { format: 'A4', orientation: 'portrait', margin: '1cm' };

    page.open(address, function (status) {
        if (status !== 'success') {
            console.log('Unable to load the address!');
            phantom.exit();
        } else {
            window.setTimeout(function () {
                page.render(output, {format: "pdf"});
                phantom.exit();
            }, 200);
        }
    });
}
