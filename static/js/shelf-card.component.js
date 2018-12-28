Vue.component('shelf-card', {
  props: ['lastView'],
  data: function() {
    return {
        books: {},
    }
  },
  mounted: function() {
    this.fetchBooksInfo();
    /*
    $(".sliding-link").on('click',function(e) {
      this.slideLink(e);
    });
    */
  },
  beforeDestroy(){
    //$('.sliding-link').off('click', function (e) {});
  },
  methods: {
    slideLink: function(tag) {
      console.log("tag:",tag);
      var aid = $(".card[jumptag='"+tag+"']");
      $('html,body').animate({scrollTop: $(aid).offset().top},'slow');
    },
    loadBook: function(book) {
      var self = this;
      console.log('inside shelf-comp.loadBook:',book)
      self.$emit('load-book', book);
    },
    fetchBooksInfo: function(event) {
      var self = this ;
      console.log("Inside fetchBooksInfo()");
      if (localStorage.getItem("books")) {
        try {
          console.log("load books from cache")
          self.books = JSON.parse(localStorage.getItem('books'));
          self.$nextTick(function() {
            console.log("finish loading books from cache");
            if (self.lastView) {
              console.log("jump to last view book:", self.lastView);
              var aid = $("a#"+self.lastView).parent("div");
              if (aid)
                $('html,body').animate({scrollTop: $(aid).offset().top},'slow');        
            }
          });
        } catch(e) {
          console.log("load cache ex:", e);
          localStorage.removeItem('books');
          self.books = null;
        }
      } else {
        console.log("unable to load books from cache")
        self.books = null;
      }
      if (!self.books) {
        console.log("load chapter from network");
        fetch('/api/books', {
          method: 'GET'//,
          /*
          headers: {
            "Content-Type": "application/json",
          }*/
        }).then(
          function(response) {
            if (response.status !== 200) {
              console.log('Looks like there was a problem. Status Code: ' + response.status);
              return;
            }
            // Examine the text in the response
            response.json().then(function(data) {
              //console.log(data);
              self.books = data;
              const parsed = JSON.stringify(data);
              localStorage.setItem("books", parsed);

              self.$nextTick(function() {
                console.log("finish loading books from network");
                if (self.lastView) {
                  console.log("jump to last view book:", self.lastView);
                  var aid = $("a#"+self.lastView).parent("div");
                  if (aid)
                    $('html,body').animate({scrollTop: $(aid).offset().top},'slow');
                }
              });
            });
          }
        ).catch(function(err) {
          console.log('Fetch Error :-S', err);
        });
      }
    }
  },
  filters: {
    firstChar: function(s) {
      return s.toUpperCase()[0];
    }
  },
  computed: {
    azIndex: function() {
      var idx = [];
      if (this.books) {
        var bL = this.books.length;
        for (var i=0;i<bL;i++) {
          var fC = this.books[i].toUpperCase()[0];
          if (idx.indexOf(fC) === -1) idx.push(fC);
        }
        return idx;  
      }
      return [];
    }
  },
  template: `<div class="main-template">
              <div class="row" id="tableOfContent">
                <a href="#" class="sliding-link" v-for="(idx) in azIndex" v-bind:jumptag="idx" @click="slideLink(idx)">{{idx}}</a>
              </div>
              <div class="row">
              <div class="col-sm-4" v-for="(item) in books"  >
              <div class="card" v-bind:jumptag="item | firstChar" >
                <div class="card-title card-header text-monospace text-truncate d-flex" >
                  <a href="#tableOfContent" class="link-top"><i class="material-icons">arrow_upward</i></a>
                  <a v-bind:id="item" class="flex-grow-1">{{item}}</a>
                </div>
                <div class="card-body">
                <a href="#" @click="loadBook(item)"><img class="card-img-bottom" :src="'/pages/'+item+'/'+item+'.thumb'"/></a>
                </div>
              </div>
              </div>
              </div>    
            </div>
          `
})


