Vue.component('book-view', {
  props: ['title'],
  data: function() {
    return {
        activeClass: 'active',
        bookTitle: '',
        chapters: [],
        activeChapter: '',
        viewingPages:[],
        activePage:0,
        jumpToPage:0
    }
  },
  mounted: function() {
    
    $('#dismiss, .overlay, .gallery').on('click', function () {
      $('#sidebar').removeClass('active');
      $('.overlay').removeClass('active');
    });

    $('#sidebarCollapse').on('click', function () {
        $('#sidebar').addClass('active');
        $('.overlay').addClass('active');
        $('.collapse.in').toggleClass('in');
        $('a[aria-expanded=true]').attr('aria-expanded', 'false');
    });
    $('#carouselImg').on('slid.bs.carousel', this.onSlid);
    
    var hasBookmark = this.loadFromBookmark();
    if (!hasBookmark)
      this.fetchChapterInfo();
  },
  beforeDestroy(){
    $('#dismiss, .overlay, .gallery').off('click', function () {
      $('#sidebar').removeClass('active');
      $('.overlay').removeClass('active');
    });

    $('#sidebarCollapse').off('click', function () {
        $('#sidebar').addClass('active');
        $('.overlay').addClass('active');
        $('.collapse.in').toggleClass('in');
        $('a[aria-expanded=true]').attr('aria-expanded', 'false');
    });
    $('#carouselImg').off('slid.bs.carousel', this.onSlid);
  },
  methods: {
    loadFromBookmark: function() {
      if (localStorage.getItem("bookmark-"+this.title)) {
        try {
          console.log("load bookmark from cache")
          var bookmark = JSON.parse(localStorage.getItem("bookmark-"+this.title));
          console.log(bookmark);
          //self.activeChapter = bookmark["active_chapter"];
          //self.activePage = bookmark["active_page"];
          this.fetchChapterInfo(bookmark["active_chapter"]);
          this.jumpToPage = bookmark["active_page"];
          //this.fetchPages(bookmark["active_chapter"],bookmark["active_page"]);
          return true;
        } catch(e) {
          console.log("load bookmark from cache ex:", e);
          localStorage.removeItem("bookmark-"+this.title);
        }
      } else {
        console.log("unable to load from cache")
      }
      return false;
    },
    onSlid(event) {
      //console.log("slid event", event);
      this.activePage = event.to ;
    },
    onSlide() {
      console.log("slide event");
    },
    loadShelf: function() {
      var self = this;
      console.log('inside book-view.loadShelf');
      self.$emit('load-shelf');
    },
    fetchChapterInfo: function(lastChapter) {
      var self = this ;
      console.log("Inside fetchChapterInfo():", self.title);
      console.log("fetchurl:", '/api/books/'+self.title)
      console.log("lastChapter:", lastChapter);
      if (localStorage.getItem('chapters-'+self.title)) {
        try {
          console.log("load chapter from cache")
          self.chapters = JSON.parse(localStorage.getItem('chapters-'+self.title));
          if (lastChapter) {
            self.activeChapter = lastChapter;
          } else {
            self.activeChapter = self.chapters[0];
          }
        } catch(e) {
          console.log("load cache ex:", e);
          localStorage.removeItem('chapters-'+self.title);
          self.chapters = null;
        }
      } else {
        console.log("unable to load from cache")
        self.chapters = null;
      }
      if (!self.chapters) {
        console.log("load chapter from network")
        fetch('/api/books/'+self.title, {
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
                self.chapters = data;
                if (lastChapter) {
                  self.activeChapter = lastChapter;
                } else {
                  self.activeChapter = data[0];
                }
                const parsed = JSON.stringify(data);
                localStorage.setItem("chapters-"+self.title, parsed);
              });
            }
          )
          .catch(function(err) {
            console.log('Fetch Error :-S', err);
          });
      }
    },
    prevChapter: function() {
      var self = this ;
      console.log("Inside prevChapter(): currentChapter:", self.activeChapter);
      var currentIdx = self.chapters.indexOf(self.activeChapter);
      console.log("currentIdx:", currentIdx);

      if (currentIdx == 0) {
        console.log("First chapter already");
      } else {
        self.activeChapter = self.chapters[currentIdx-1];
      }
    },
    nextChapter: function() {
      var self = this ;
      console.log("Inside nextChapter(): currentChapter:", self.activeChapter);
      var currentIdx = self.chapters.indexOf(self.activeChapter);
      console.log("currentIdx:", currentIdx);

      if (currentIdx == self.chapters.length-1) {
        console.log("Last chapter already");
      } else {
        self.activeChapter = self.chapters[currentIdx+1];
      }
    },
    gotoChapter: function(chapter) {
      var self = this ;
      console.log("Inside gotoChapter(): currentChapter:", self.activeChapter);
      console.log("goto:", chapter);

      if (self.chapters) {
        var found = self.chapters.indexOf(chapter);
        if (found >=0) {
          self.activeChapter = chapter;
        }
      }
    },
    fetchPages: function(chapter, lastPage) {
      var self = this ;
      console.log("Inside fetchPage():", chapter, self.title);
      console.log("fetchurl:", '/api/books/'+self.title+"/"+chapter);
      console.log("lastPage:", lastPage);
      fetch('/api/books/'+self.title+"/"+chapter, {
        method: 'GET'//,
      }).then(
          function(response) {
            if (response.status !== 200) {
              console.log('Looks like there was a problem. Status Code: ' + response.status);
              return;
            }
            // Examine the text in the response
            response.json().then(function(data) {
              //console.log(data);
              self.viewingPages = data;
              //self.activeChapter = chapter ;
              if (lastPage) {
                self.activePage = lastPage;
                $('#carouselImg div.carousel-item:nth-child('+lastPage+')').addClass('active');
              } else {
                self.activePage = 0;
                $('#carouselImg div.carousel-item:first').addClass('active');
              }
            });
          }
        )
        .catch(function(err) {
          console.log('Fetch Error :-S', err);
        });
    }
  },
  computed: {
    totalPages: function() {
      return this.viewingPages.length;
    },
    currentChapter: function() {
      if (this.chapters) {
        return this.chapters.indexOf(this.activeChapter);
      }
      return 0;
    }
  },
  filters: {
  },
  watch: {
    activeChapter: function(newChapter, oldChapter) {
      console.log("watch change chapters:",oldChapter,"->",newChapter);
      console.log("jumpPage:", this.jumpToPage);
      this.fetchPages(newChapter, this.jumpToPage);
      this.jumpToPage = 0;
    },
    activePage: function(newPage, oldPage) {
      console.log("turn page:",oldPage,"->",newPage);
      var bookmark={"active_chapter":this.activeChapter, "active_page":this.activePage};
      var parsed = JSON.stringify(bookmark);
      localStorage.setItem("bookmark-"+this.title,parsed);
      window.scrollTo({ top: 0, behavior: 'smooth' });
    }
  },
  template: `
  <div>
  <!-- Sidebar  -->
  <nav id="sidebar">
      <div id="dismiss">
          <i class="material-icons">arrow_left</i>
      </div>

      <div class="sidebar-header">
          <h4>Chapters</h4>
      </div>

      <ul class="list-unstyled components">
          <!--
          <li class="active">
              <a href="#homeSubmenu" data-toggle="collapse" aria-expanded="false">Home</a>
              <ul class="collapse list-unstyled" id="homeSubmenu">
                  <li>
                      <a href="#">Home 1</a>
                  </li>
                  <li>
                      <a href="#">Home 2</a>
                  </li>
              </ul>
          </li>
          -->
          <li v-for="item in chapters" v-bind:class="item==activeChapter?activeClass:''"><a href="#" @click="gotoChapter(item)">{{item}}</a></li>
      </ul>
  </nav>
  <div id="content">
    <div class="container sticky-top ">
    <nav class="navbar navbar-light bg-light navbar-expand-sm d-flex align-items-baseline pl-0">
      <a href="#" @click="loadShelf()"><i class="material-icons">dashboard</i></a>
      <a href="#" @click="prevChapter()"  class="pl-2"><i class="material-icons pr-2">arrow_back</i></a>
      <div class="flex-grow-1 text-center small">
        <a href="#" id="sidebarCollapse" class="" >{{activePage+1}}/{{totalPages}} C{{currentChapter+1}} {{title}}</a>
      </div>
      <a href="#" @click="nextChapter()"><i class="material-icons pl-2">arrow_forward</i></a>
    </nav>
    </div>
    <h2></h2>
    <div class="gallery">
      <!--
      <ul>
      <li v-for="img in viewingPages"><img :src="'/pages/'+img"/></li>
      </ul>
      -->
      <div id="carouselImg" class="carousel slide" data-ride="carousel" data-wrap="false" data-interval="10000">
        <div class="carousel-inner">
          <div class="carousel-item" v-for="(img,idx) in viewingPages" :key="idx" v-bind:class="idx==activePage?activeClass:''">
            <img class="d-block w-100 img-fluid" :src="'/pages/'+img"/>
          </div>
        </div>
        <a class="carousel-control-prev" href="#carouselImg" role="button" data-slide="prev">
          <span class="carousel-control-prev-icon" aria-hidden="true"></span>
          <span class="sr-only">Previous</span>
        </a>
        <a class="carousel-control-next" href="#carouselImg" role="button" data-slide="next">
          <span class="carousel-control-next-icon" aria-hidden="true"></span>
          <span class="sr-only">Next</span>
        </a>
      </div>
    </div>
    </div>
    <div class="overlay"></div>
  </div>
  `
})


