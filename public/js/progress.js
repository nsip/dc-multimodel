// progress.js


// 
// fill the progress data table
// 
function fillProgressTable(tgroup_id) {

    var $progress = $("#progress-table");
    $($progress).empty();

    $.ajax({
        method: "POST",
        url: "http://localhost:1340/graphql",
        contentType: "application/json",
        data: JSON.stringify({
            query: progressQuery,
            variables: { "teaching_group_id": tgroup_id }
        })
    }).done(function(result) {
        // fill teaching group selector
        // console.log(result);
        $.each(result.data.reports, function() {
            // console.log(this);
            // add a title - the teaching group
            $($progress).append("<h5>" + this.teaching_group_name + "</h5>");
            // iterate the progress results
            $.each(this.assignment_results, function() {
                // add a title - the assigment name
                $($progress).append("</br>");
                $($progress).append("<h5>Assignment Name:  " + this.assignment_name + "</h5>");
                // add a results table
                var $table = $('<table>').addClass('highlight showmodal');
                // add table header
                $table.append("<thead><tr><th>Name</th><th>Success</th><th>Completed</th><th>Score</th><th>Time On Task</th><th>Days Absent</th></tr></thead>")

                // iterate the student results for this assignment
                var $studentResults = $('<tbody>');
                // $table.append($studentResults);
                $.each(this.student_results, function() {
                    var $resultsRow = $('<tr>');
                    $resultsRow.append("<td>" + this.student_name + "</td>");
                    $resultsRow.append("<td>" + this.result_event.result.success + "</td>");
                    $resultsRow.append("<td>" + this.result_event.result.completion + "</td>");
                    $resultsRow.append("<td>" + this.result_event.result.score.scaled + "</td>");
                    $resultsRow.append("<td>" + minutesFromIsoDuration(this.result_event.result.duration) + " minutes</td>");
                    $resultsRow.append("<td>" + this.absence_days + "</td>");
                    $studentResults.append($resultsRow);
                    $table.append($studentResults);
                });
                $progress.append($table);
                // console.log($table);
            });
        });
    });
    initProgressTableHandler();

}


// 
// helpers for parsing ISO durations from xAPI
// 
var regex = /P((([0-9]*\.?[0-9]*)Y)?(([0-9]*\.?[0-9]*)M)?(([0-9]*\.?[0-9]*)W)?(([0-9]*\.?[0-9]*)D)?)?(T(([0-9]*\.?[0-9]*)H)?(([0-9]*\.?[0-9]*)M)?(([0-9]*\.?[0-9]*)S)?)?/

minutesFromIsoDuration = function(duration) {
    var matches = duration.match(regex);

    return parseFloat(matches[14]) || 0;
}


// 
// the GQL query
// 
var progressQuery = `
query progressReport($teaching_group_id: String!) {
  reports(teaching_group_id: $teaching_group_id) {
    teaching_group_name
    assignment_results {
      assignment_name
      student_results {
        student_name
        absence_days
        result_event{
         result {
          completion
          success
          duration
          score {
            scaled
            # min
            # max
          }
        }
        }
      }
    }
  }
}
`