// teachers.js

// 
// fill the teachers selector
// 
function fillTeacherSelector() {
    var $dropdown = $("#teacher-select");
    $($dropdown).empty();

    $.ajax({
        method: "POST",
        url: "http://localhost:1340/graphql",
        contentType: "application/json",
        data: JSON.stringify({
            query: teacherQuery
        })
    }).done(function(result) {
        // fill the teachers dropdown
        $.each(result.data.teachers, function() {
            // console.log(this);
            $dropdown.prepend($("<option />").val(this.RefId).text(this.PersonInfo.Name.GivenName +
                " " + this.PersonInfo.Name.FamilyName));
        });
        // have to re-initialise component to render
        $($dropdown).formSelect();
    });
}


var teacherQuery = `
			query teachersQuery {
			  teachers {
			    RefId
			    # LocalId
			    # EmploymentStatus
			    PersonInfo{
			      Name{
			        FamilyName
			        GivenName
			      }
			    }
			  }
			}`
