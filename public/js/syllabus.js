// syllabus.js

// 
// fill syllabus view
// 
function fillSyllabusView(teacher_guid) {

    var $concepts = $("#concepts");
    $($concepts).empty();

    var $skills = $("#skills");
    $($skills).empty();

    var $tools = $("#tools");
    $($tools).empty();

    $("#status-message").empty();

    // only hsie syllabus supported
    var $selected = $("#teaching-group-select option:selected").text();
    var $learning_area = "";
    if ($selected.indexOf("Geography") != -1) {
        $learning_area = "geography";
    }
    if ($selected.indexOf("History") != -1) {
        $learning_area = "history";
    }

    if (($learning_area != "history") && ($learning_area != "geography")) {
        $("#status-message").append("<p>Syllabus tracking not available for this Learning Area</p>");
        return
    }


    $.ajax({
        method: "POST",
        url: "http://localhost:1330/graphql",
        contentType: "application/json",
        data: JSON.stringify({
            query: syllabusQuery,
            variables: {
                "state": "nsw",
                "learning_area": "hsie",
                "subject": $learning_area,
                "stage": "4"
            }
        })
    }).done(function(result) {
        // fill popoup view
        // console.log(result);
        // add column header
        $concepts.append($("<p><h6>Concepts</h6></p>"));
        $.each(result.data.overview.concepts, function() {
            // console.log(this);
            var $conceptsBlock = $('<div>');
            // $conceptsBlock.append($("<label><input type='checkbox' class='filled-in' checked='checked' /><span>" +
            // this.name + "</span></label>"));
            $conceptsBlock.append(makeCheckbox(this.name));
            var $list = $('<ul>').addClass("collection");
            var $listItem = $("<li>").addClass("collection-item");
            $listItem.text(this.description);
            $($list).append($listItem);
            $conceptsBlock.append($list);
            $concepts.append($($conceptsBlock));
        });
        // add column header
        $skills.append($("<p><h6>Inquiry Skills</h6></p>"));
        $.each(result.data.overview.inquiry_skills, function() {
            var $skillsBlock = $('<div>');
            // console.log(this);
            // $skillsBlock.append($("<label><input type='checkbox' class='filled-in' checked='checked' /><span>" +
            // this.name + "</span></label>"));
            $skillsBlock.append(makeCheckbox(this.name));
            var $list = $('<ul>').addClass("collection");
            $.each(this.skills, function() {
                var $listItem = $("<li>").addClass("collection-item");
                $listItem.text(this.skill);
                $($list).append($listItem);
            });
            $skillsBlock.append($list);
            $skills.append($skillsBlock)
        });
        // add column header
        $tools.append($("<p><h6>Tools</h6></p>"));
        $.each(result.data.overview.tools, function() {
            var $toolsBlock = $('<div>');
            // console.log(this);
            // $toolsBlock.append($("<label><input type='checkbox' class='filled-in' checked='checked' /><span>" +
            // this.name + "</span></label>"));
            $toolsBlock.append(makeCheckbox(this.name));
            var $list = $('<ul>').addClass("collection");
            if ($learning_area == "geography") { //history returns null objects
                $.each(this.examples, function() {
                    // console.log(this);
                    var $listItem = $("<li>").addClass("collection-item");
                    $listItem.text(this);
                    $($list).append($listItem);
                });
                $toolsBlock.append($list);
                $tools.append($toolsBlock)
            }
        });


    });

}


function makeCheckbox(labelText) {
    var $label = $("<label>");
    var $input = $("<input type='checkbox'>").addClass('filled-in');
    if (toggle()) {
        $input = $("<input type='checkbox' checked='checked'>").addClass('filled-in');
    }
    var $span = $("<span>").text(labelText);
    $($label).append($input).append($span);
    // console.log($label);
    return $label;
}

function toggle() {
    return Math.random() >= 0.5;
}


var syllabusQuery = `query OverviewQuery($state: String!, $learning_area: String!, $subject: String!, $stage: String!) {
  overview(state: $state, learning_area: $learning_area, subject: $subject, stage: $stage) {
    concepts {
      name
      description
    }
    inquiry_skills {
      name
      skills {
        skill
      }
    }
    tools {
      name
      examples
    }
    }
}`