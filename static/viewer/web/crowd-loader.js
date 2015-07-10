var crowdTakeOver = function() {
    return {
        start: function (app) {
            app.annotations.load({uri: window.location.href});
            //app.notify("SwartzNotes activated. Select text to begin annotating!");
        }
    };
};

var getQueryString = function(key) {
    key = key.replace(/[*+?^$.\[\]{}()|\\\/]/g, "\\$&"); // escape RegEx meta chars
    var match = location.search.match(new RegExp("[?&]"+key+"=([^&]+)(&|$)"));
    return match && decodeURIComponent(match[1].replace(/\+/g, " "));
};

var pageUri = function () {
    return {
        beforeAnnotationCreated: function (ann) {
            ann.uri = window.location.href;
        }
    };
};

var initializeAnnotator = function() {
	if (typeof annotator === 'undefined') {
	  /*
	  alert("Oops! it looks like you haven't built Annotator. " +
	        "Either download a tagged release from GitHub, or build the " +
	        "package by running `make`");
	  */
	  console.log("I could not load the annotator plugin for some reason. This page is probably not supported, sorry!");
	} else {
	  var app = new annotator.App();

	  app.include(crowdTakeOver);
	  app.include(pageUri);
	  app.include(annotator.storage.http, {
	    prefix: 'http://52.5.78.150/swnotes-store',
	    localSuggestionPrefix: '',
	    localSuggestionURL: getQueryString("facts"),
	  });
	  app.include(annotator.ui.main, {
	    viewerExtensions: [annotator.ui.tags.viewerExtension, annotator.ui.crowd.viewerExtension],
	    editorExtensions: [annotator.ui.tags.editorExtension, annotator.ui.crowd.editorExtension]
	  });

	  app.start();
	}
};

var getRandomColor = function() {
    var letters = '0123456789ABCDEF'.split('');
    var color = '#';
    for (var i = 0; i < 6; i++ ) {
        color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
};

/*

// This package depends on jQuery
var FactFactory = function () {
	//console.log("Creating a new FactFactory");
	this.facts = [];
	//console.log(this);
};



// getEditDistance - Copyright (c) 2011 Andrei Mackenzie

// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

 
// Compute the edit distance between the two given strings
FactFactory.getEditDistance = function(afull, bfull){
  var a = afull.split(" ");
  var b = bfull.split(" ");

  if(a.length == 0) return b.length; 
  if(b.length == 0) return a.length; 
 
  var matrix = [];
 
  // increment along the first column of each row
  var i;
  for(i = 0; i <= b.length; i++){
    matrix[i] = [i];
  }
 
  // increment each column in the first row
  var j;
  for(j = 0; j <= a.length; j++){
    matrix[0][j] = j;
  }
 
  // Fill in the rest of the matrix
  for(i = 1; i <= b.length; i++){
    for(j = 1; j <= a.length; j++){
      if(b[i-1] == a[j-1]){
        matrix[i][j] = matrix[i-1][j-1] + 2;
      } else {
        matrix[i][j] = Math.max(matrix[i-1][j-1] - 1, // substitution
                                Math.max(matrix[i][j-1], // insertion
                                         matrix[i-1][j])); // deletion
      }
    }
  }
 
  return matrix[b.length][a.length];
};


// LoadFactsFromURL loads the facts file from a specificed URL
// returns a HTML5 Promises/A+ object
 
FactFactory.prototype.LoadFactsFromURL = function(url) {
  // Return a new promise.
  var ff = this;
  return new Promise(function(resolve, reject) {
	$.getJSON(url).success(function(data) {

		data.forEach ( function(element, index, array) {
			if (element.type === "SectionTitle") {
				
			} else if (element.type === "Sentence") {
				ff.facts.push({"sentence": element.sentence, "facts": element.fact});
			}

		});

		//console.log("Resolving the Promise for LoadFactsFromURL");

		resolve(ff.facts);

	}).fail(function() {
		reject(Error("Cannot load this file!"));
	});
  });
};

// Tries its best to create all the annotations
// based on the facts it is given
FactFactory.prototype.CreateAnnotations = function(facts) {
	var potentialElements = $('#viewer').children().children().children().toArray();
	//console.log(potentialElements);
	facts.forEach ( function(element, index, array) {
		//console.log("Looking for match for sentence: "+ element.sentence);

		var bestSentence = null;
		var bestSentenceScore = 0;

		potentialElements.forEach( function(element2, index2, array2) {
			var matchScore = FactFactory.getEditDistance( $(element2).text(), element.sentence );
			if (bestSentenceScore == 0 || matchScore > bestSentenceScore) {
				bestSentenceScore = matchScore;
				bestSentence = element2;
			}
		});

		//console.log( "The best sentence is -- " + $(bestSentence).text() );
		
		$(bestSentence).css("background-color", getRandomColor());

	} );
};

myFactor = new FactFactory();

var loadFactsFile = function() {
	if (getQueryString("facts") === "") {
		console.log("You haven't sepcified a facts file to load. We're not going to make any annotations!");
		return;
	}

	myFactor.LoadFactsFromURL( getQueryString("facts") ).then( myFactor.CreateAnnotations );


	// TODO: Remove this. This is beyond horrible.
	// $('#viewer').children().children().children().each(function(index) {
	//	console.log($(this).text() + "///");
	//	$(this).css("background-color", getRandomColor());
	// });

		//console.log(data);


};

*/

// Huan: This is a horrible hack to make the annotator hooks
// because we can't be sure that the PDF is done rendering
// before the annotation overlay is done loading
// TODO: Catch this programmatically intead of using a timer.

setTimeout(function () { initializeAnnotator(); }, 3000);
//setTimeout(function () { loadFactsFile(); }, 3000);
