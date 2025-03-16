import { copyLibFiles } from '@qwik.dev/partytown/utils'; // ESM
// const { copyLibFiles } = require('@builder.io/partytown/utils'); // CommonJS

async function myBuildTask() {
  await copyLibFiles('./cmd/conneroh/_static/dist/');
}

// call the buildTask function
myBuildTask();
