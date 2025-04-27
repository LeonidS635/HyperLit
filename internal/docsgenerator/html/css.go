package html

const css = `<style>
	body {
		margin: 0;
		font-family: Arial, sans-serif;
	}
	.container {
		display: flex;
		height: 100vh;
	}
	.tree {
		width: 33%;
		padding-right: 20px;
		border-right: 2px solid #ccc;
		overflow-x: auto;
		overflow-y: auto;
	}
	.content {
		width: 67%;
		padding-left: 20px;
		overflow-y: auto;
	}
	.folder, .file {
		margin-top: 5px;
		cursor: pointer;
		display: block;
	}
	.folder {
		font-weight: bold;
	}
	.file {
		font-weight: normal;
	}
	ul {
		list-style-type: none;
		padding-left: 0;
	}
	.hidden {
		display: none;
	}
	.folder-description, .file-description {
		margin-top: 20px;
		color: #555;
	}
	.nested {
		margin-left: 20px; /* Каждый уровень вложенности увеличивает отступ */
	}
</style>`
