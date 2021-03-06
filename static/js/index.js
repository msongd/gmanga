
var app = new Vue({ 
    el: '#app',
    data: function() {
      return {
        activeView: 'shelf',
        viewingBook: ''
      }
    },
    mounted: function() {
    },
    methods: {
        loadBook: function(book) {
            var self = this;
            console.log("Inside root.loadBook:", book);
            self.activeView = 'book';
            self.viewingBook = book;
        },
        loadShelf: function(ev) {
            var self = this;
            console.log("Inside root.loadShelf:", ev);
            self.activeView = 'shelf';
            self.viewingBook = ev;
        }
    },
});
