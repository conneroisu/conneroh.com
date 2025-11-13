#!/usr/bin/env bun

/**
 * Sanity Connection Verification Script
 *
 * This script verifies that:
 * 1. Environment variables are properly configured
 * 2. API token is valid and has read access
 * 3. We can successfully fetch data from Sanity CMS
 * 4. Reports available content (posts, projects, tags, employments)
 */

import { fetchSanityContent } from '../src/lib/sanity-server';
import { listPosts, listProjects, listTags, listEmployments, countPosts, countProjects } from '../src/lib/sanity-queries';
import type { Post, Project, Tag, Employment } from '../src/lib/sanity-types';

// ANSI color codes for terminal output
const colors = {
  reset: '\x1b[0m',
  green: '\x1b[32m',
  red: '\x1b[31m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
  cyan: '\x1b[36m',
  bold: '\x1b[1m',
};

function logSuccess(message: string) {
  console.log(`${colors.green}✓${colors.reset} ${message}`);
}

function logError(message: string) {
  console.log(`${colors.red}✗${colors.reset} ${message}`);
}

function logInfo(message: string) {
  console.log(`${colors.blue}ℹ${colors.reset} ${message}`);
}

function logSection(title: string) {
  console.log(`\n${colors.bold}${colors.cyan}${title}${colors.reset}`);
  console.log('─'.repeat(title.length));
}

async function verifyEnvironmentVariables(): Promise<boolean> {
  logSection('Step 1: Verifying Environment Variables');

  const projectId = process.env.VITE_SANITY_PROJECT_ID;
  const dataset = process.env.VITE_SANITY_DATASET;
  const apiToken = process.env.SANITY_API_TOKEN;

  let allPresent = true;

  if (projectId) {
    logSuccess(`VITE_SANITY_PROJECT_ID: ${projectId}`);
  } else {
    logError('VITE_SANITY_PROJECT_ID is missing');
    allPresent = false;
  }

  if (dataset) {
    logSuccess(`VITE_SANITY_DATASET: ${dataset}`);
  } else {
    logError('VITE_SANITY_DATASET is missing');
    allPresent = false;
  }

  if (apiToken) {
    logSuccess(`SANITY_API_TOKEN: ${apiToken.substring(0, 20)}...${apiToken.substring(apiToken.length - 10)} (${apiToken.length} chars)`);
  } else {
    logError('SANITY_API_TOKEN is missing');
    allPresent = false;
  }

  return allPresent;
}

async function verifyConnection(): Promise<boolean> {
  logSection('Step 2: Testing API Connection');

  try {
    // Simple query to test connection
    const result = await fetchSanityContent<{ _id: string }[]>('*[_type == "post"][0..0]{_id}', {}, {
      retry: false,
      logErrors: false
    });

    logSuccess('Successfully connected to Sanity API');
    logInfo('API token is valid and has read access');
    return true;
  } catch (error) {
    logError('Failed to connect to Sanity API');
    if (error instanceof Error) {
      logError(`Error: ${error.message}`);
    }
    return false;
  }
}

async function fetchContentCounts(): Promise<{
  posts: number;
  projects: number;
  tags: number;
  employments: number;
}> {
  logSection('Step 3: Fetching Content Counts');

  const [postsCount, projectsCount, tagsResult, employmentsResult] = await Promise.all([
    fetchSanityContent<number>(countPosts(), {}, { logErrors: false }),
    fetchSanityContent<number>(countProjects(), {}, { logErrors: false }),
    fetchSanityContent<Tag[]>(listTags(), {}, { logErrors: false }),
    fetchSanityContent<Employment[]>(listEmployments(), {}, { logErrors: false }),
  ]);

  const counts = {
    posts: postsCount,
    projects: projectsCount,
    tags: tagsResult.length,
    employments: employmentsResult.length,
  };

  logSuccess(`Posts: ${counts.posts}`);
  logSuccess(`Projects: ${counts.projects}`);
  logSuccess(`Tags: ${counts.tags}`);
  logSuccess(`Employments: ${counts.employments}`);

  return counts;
}

async function fetchSampleContent(): Promise<void> {
  logSection('Step 4: Fetching Sample Content');

  // Fetch first 3 items of each type
  const [posts, projects, tags, employments] = await Promise.all([
    fetchSanityContent<Post[]>(listPosts() + '[0..2]', {}, { logErrors: false }),
    fetchSanityContent<Project[]>(listProjects() + '[0..2]', {}, { logErrors: false }),
    fetchSanityContent<Tag[]>(listTags() + '[0..2]', {}, { logErrors: false }),
    fetchSanityContent<Employment[]>(listEmployments() + '[0..2]', {}, { logErrors: false }),
  ]);

  if (posts.length > 0) {
    logInfo('\nSample Posts:');
    posts.forEach((post, i) => {
      console.log(`  ${i + 1}. ${post.title || 'Untitled'} (${post.slug?.current || 'no-slug'})`);
    });
  } else {
    logInfo('No posts found');
  }

  if (projects.length > 0) {
    logInfo('\nSample Projects:');
    projects.forEach((project, i) => {
      console.log(`  ${i + 1}. ${project.title || 'Untitled'} (${project.slug?.current || 'no-slug'})`);
    });
  } else {
    logInfo('No projects found');
  }

  if (tags.length > 0) {
    logInfo('\nSample Tags:');
    tags.forEach((tag, i) => {
      console.log(`  ${i + 1}. ${tag.title || 'Untitled'} (${tag.slug?.current || 'no-slug'})`);
    });
  } else {
    logInfo('No tags found');
  }

  if (employments.length > 0) {
    logInfo('\nSample Employments:');
    employments.forEach((employment, i) => {
      console.log(`  ${i + 1}. ${employment.title || 'Untitled'} (${employment.slug?.current || 'no-slug'})`);
    });
  } else {
    logInfo('No employments found');
  }
}

async function main() {
  console.log(`\n${colors.bold}${colors.cyan}Sanity CMS Connection Verification${colors.reset}`);
  console.log('═'.repeat(40));

  try {
    // Step 1: Verify environment variables
    const envVarsOk = await verifyEnvironmentVariables();
    if (!envVarsOk) {
      logError('\nVerification failed: Missing required environment variables');
      logInfo('Please ensure .env.local contains all required variables');
      process.exit(1);
    }

    // Step 2: Test connection
    const connectionOk = await verifyConnection();
    if (!connectionOk) {
      logError('\nVerification failed: Cannot connect to Sanity API');
      logInfo('Please check that your API token is valid and has read permissions');
      process.exit(1);
    }

    // Step 3: Fetch content counts
    const counts = await fetchContentCounts();

    // Step 4: Fetch and display sample content
    await fetchSampleContent();

    // Success summary
    logSection('Verification Complete');
    logSuccess('All checks passed!');
    logInfo('\nSanity CMS integration is ready to use');

    const totalContent = counts.posts + counts.projects + counts.tags + counts.employments;
    logInfo(`Total content items: ${totalContent}`);

    if (totalContent === 0) {
      logInfo('\n' + colors.yellow + 'Note: No content found in Sanity. You may want to add some content through the Sanity Studio.' + colors.reset);
    }

  } catch (error) {
    logError('\nUnexpected error during verification');
    if (error instanceof Error) {
      console.error(error);
    }
    process.exit(1);
  }
}

// Run the script
main();
