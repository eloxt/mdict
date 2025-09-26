package replacer

const WordDefinitionTemplate = `
<html>
<head>
<script lang="javascript">

function _play_sound(url) {
	var audio = new Audio(url);
	audio.play();
}

var __TOPFRAME_SECURE_ORIGIN__ = "*";
function _entry_jump(word, dict_id) {
	if (window.top){
		window.top.postMessage({"evtype":"_INNER_FRAME_MSG_EVTP_ENTRY_JUMP", "word":word, "dict_id":dict_id},__TOPFRAME_SECURE_ORIGIN__ )
	}
}
</script>
</head>
<body>
%s
</body>
</html>
`
