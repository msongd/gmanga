Vue.component('shelf-card', {
  data: function() {
    return {
        books: {},
    }
  },
  mounted: function() {
    this.fetchBooksInfo();
  },
  computed: {
  },
  methods: {
    loadBook: function(book) {
      var self = this;
      console.log('inside shelf-comp.loadBook:',book)
      self.$emit('load-book', book);
    },
    fetchBooksInfo: function(event) {
      var self = this ;
      console.log("Inside fetchBooksInfo()");
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
            });
          }
        )
        .catch(function(err) {
          console.log('Fetch Error :-S', err);
        });
    }
  },
  filters: {
  },
  template: `<div class="main-template">              
              <div class="row">
              <div class="col-sm-4" v-for="(item) in books"  >
              <div class="card"  >
                <div class="card-title card-header text-monospace text-truncate" >
                  <a v-bind:href="'#'+item" @click="loadBook(item)">{{item}}</a>
                </div>
                <div class="card-body">
                <a v-bind:href="'#'+item" @click="loadBook(item)"><img class="card-img-bottom" :src="'/pages/'+item+'/'+item+'.thumb'"/></a>
                </div>
              </div>
              </div>
              </div>    
            </div>
          `
})


