(function($) {
    $.addParser({
        id: 'filesize',
        re: /^(\d+(?:\.\d+)?) (bytes?|kiB|MiB|GiB|TiB|PiB)$/,
        is: function(s) {
            return this.re.test(s);
        },
        format: function(s) {
            var matches = this.re.exec(s);
            if (!matches) {
                return 0;
            }
            var size = parseFloat(matches[1]);
            var suffix = matches[2];

            switch(suffix) {
                /* intentional fall-through at each level */
                case 'PiB':
                    size *= 1024;
                case 'TiB':
                    size *= 1024;
                case 'GiB':
                    size *= 1024;
                case 'MiB':
                    size *= 1024;
                case 'kiB':
                    size *= 1024;
            }
            return size;
        },
        type: 'numeric'
    });
})($.tablesorter);

