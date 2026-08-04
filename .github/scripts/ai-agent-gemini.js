const fs = require('fs');
const path = require('path');
const { GoogleGenAI, Type } = require('@google/genai');

const ai = new GoogleGenAI({ apiKey: process.env.GEMINI_API_KEY });

// Ignore the following paths when scanning the project context
const IGNORE_PATHS = ['node_modules', '.git', '.github', 'dist', 'build', 'package-lock.json'];


function getProjectContext(dir = '.') {
  let context = '';
  const entries = fs.readdirSync(dir, { withFileTypes: true });

  for (const entry of entries) {
    if (IGNORE_PATHS.includes(entry.name)) continue;

    const fullPath = path.join(dir, entry.name);

    if (entry.isDirectory()) {
      context += getProjectContext(fullPath);
    } 
    else if (entry.isFile()) {
      try {
        const content = fs.readFileSync(fullPath, 'utf8');
        context += `\n--- START OF FILE: ${fullPath} ---\n`;
        context += content;
        context += `\n--- END OF FILE: ${fullPath} ---\n`;
      } catch (e) {
        // Ignore files that can't be read (like binaries)
      }
    }
  }
  return context;
}

async function run() {
  const title = process.env.ISSUE_TITLE;
  const body = process.env.ISSUE_BODY;

  console.log('Scanning project context...');
  const projectContext = getProjectContext('.');

  const prompt = `
  You are an expert automated software developer.
  
  The current structure and content of the repository is shown below:
  ${projectContext}

  ---
  TASK:
  Modify or create the necessary files to resolve the following GitHub issue.

  Issue title: ${title}
  Issue description: ${body}

  Make sure to maintain the project's coding style and reuse existing functions where necessary.
  `;

  console.log('Sending prompt to Gemini...');

  // JSON schema for the expected response
  const response = await ai.models.generateContent({
    model: 'gemini-2.5-flash',
    contents: prompt,
    config: {
      responseMimeType: 'application/json',
      responseSchema: {
        type: Type.OBJECT,
        properties: {
          files: {
            type: Type.ARRAY,
            description: 'List of created or modified files',
            items: {
              type: Type.OBJECT,
              properties: {
                path: {
                  type: Type.STRING,
                  description: 'Relative file path (e.g., src/index.js)'
                },
                content: {
                  type: Type.STRING,
                  description: 'Full content of the generated or modified file'
                }
              },
              required: ['path', 'content']
            }
          }
        },
        required: ['files']
      }
    }
  });

  const result = JSON.parse(response.text);

  if (!result.files || result.files.length === 0) {
    console.log('Gemini did not propose any changes to the files.');
    return;
  }

  for (const file of result.files) {
    const dir = path.dirname(file.path);
    if (!fs.existsSync(dir)) {
      fs.mkdirSync(dir, { recursive: true });
    }
    fs.writeFileSync(file.path, file.content, 'utf8');
    console.log(`File updated/created: ${file.path}`);
  }
}

run().catch((err) => {
  console.error('Error during the agent execution:', err);
  process.exit(1);
});