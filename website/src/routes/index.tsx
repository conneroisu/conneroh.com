import { createFileRoute, Link } from '@tanstack/solid-router'
import { For, Show, Suspense } from 'solid-js'
import { fetchSanityContent } from '@/lib/sanity-server'
import { listProjects, listPosts, listTags, listEmployments } from '@/lib/sanity-queries'
import type { Project, Post, Tag, Employment } from '@/lib/sanity-types'
import ProjectCard from '@/components/ProjectCard'
import PostCard from '@/components/PostCard'
import TagCard from '@/components/TagCard'
import EmploymentTimeline from '@/components/EmploymentTimeline'
import ContactForm from '@/components/ContactForm'

export const Route = createFileRoute('/')({
  component: HomePage,
  loader: async () => {
    // Fetch all data in parallel
    const [projects, posts, tags, employments] = await Promise.all([
      fetchSanityContent<Project[]>(listProjects()),
      fetchSanityContent<Post[]>(listPosts()),
      fetchSanityContent<Tag[]>(listTags()),
      fetchSanityContent<Employment[]>(listEmployments()),
    ])

    return {
      projects: projects.slice(0, 6), // First 6 projects
      posts: posts.slice(0, 4), // First 4 posts
      tags: tags.slice(0, 28), // First 28 tags
      employments,
    }
  },
})

function HomePage() {
  const data = Route.useLoaderData()

  return (
    <div class="min-h-screen bg-gray-900">
      {/* Hero Section */}
      <section class="bg-gradient-to-b from-gray-900 to-gray-800 py-20">
        <div class="container mx-auto px-4">
          <div class="flex flex-col lg:flex-row items-center justify-between gap-12">
            {/* Left side - Text content */}
            <div class="lg:w-1/2 space-y-6">
              <h1 class="text-4xl md:text-5xl lg:text-6xl font-bold text-white">
                Conner Ohnesorge
              </h1>
              <p class="text-xl text-gray-300">
                Electrical Engineer & Software Developer
              </p>
              <p class="text-xl text-gray-300">
                Electrical Engineering Bachelors Degree and Minor in Computer Science from Iowa State University
              </p>
              <div class="flex gap-4 pt-4">
                <Link
                  to="/projects"
                  class="bg-white text-gray-900 px-6 py-3 rounded-md font-semibold hover:bg-gray-100 transition-colors"
                >
                  View Projects
                </Link>
                <a
                  href="#contact"
                  class="border border-gray-600 text-gray-300 px-6 py-3 rounded-md font-semibold hover:bg-gray-800 transition-colors"
                >
                  Contact Me
                </a>
              </div>
            </div>

            {/* Right side - Profile image */}
            <div class="lg:w-1/2 flex justify-center relative">
              <div class="relative">
                {/* Decorative background circles */}
                <div class="absolute -top-4 -right-4 w-32 h-32 bg-purple-600 rounded-full opacity-20 blur-2xl"></div>
                <div class="absolute -bottom-4 -left-4 w-32 h-32 bg-pink-600 rounded-full opacity-20 blur-2xl"></div>

                {/* Profile image */}
                <img
                  src="/dist/hero.jpeg"
                  alt="Conner Ohnesorge"
                  class="h-64 w-64 md:h-80 md:w-80 rounded-full object-cover border-4 border-gray-700 relative z-10"
                />
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Featured Projects Section */}
      <section class="bg-gray-800 py-16">
        <div class="container mx-auto px-4">
          <div class="mb-12 flex justify-between items-end">
            <div>
              <h2 class="text-3xl font-bold text-white mb-4">Featured Projects</h2>
              <p class="text-gray-300">Showcase of my recent work and contributions</p>
            </div>
            <Link
              to="/projects"
              class="text-emerald-400 hover:text-emerald-300 transition-colors flex items-center gap-2"
            >
              View All Projects
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
              </svg>
            </Link>
          </div>

          <Suspense fallback={<div class="text-white text-center">Loading projects...</div>}>
            <Show
              when={data().projects && data().projects.length > 0}
              fallback={<div class="text-gray-400 text-center">No projects available</div>}
            >
              <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
                <For each={data().projects}>
                  {(project) => (
                    <ProjectCard
                      title={project.title || 'Untitled Project'}
                      description={project.description}
                      bannerPath={project.banner?.asset?._ref ? `https://cdn.sanity.io/images/${project.banner.asset._ref}` : undefined}
                      tags={project.tags}
                      href={`/projects/${project.slug?.current}`}
                    />
                  )}
                </For>
              </div>
            </Show>
          </Suspense>
        </div>
      </section>

      {/* Employment Timeline Section */}
      <section class="bg-gray-900 py-16">
        <div class="container mx-auto px-4">
          <div class="mb-12">
            <h2 class="text-3xl font-bold text-white mb-4">Professional Experience</h2>
            <p class="text-gray-300">My journey through various roles in engineering and technology</p>
          </div>

          <Suspense fallback={<div class="text-white text-center">Loading experience...</div>}>
            <Show
              when={data().employments && data().employments.length > 0}
              fallback={<div class="text-gray-400 text-center">No employment records available</div>}
            >
              <EmploymentTimeline
                items={data().employments.map((emp) => ({
                  title: emp.title || 'Untitled Position',
                  description: emp.description,
                  startDate: (emp as any).startDate || emp.createdAt || new Date().toISOString(),
                  endDate: emp.endDate,
                  tags: emp.tags,
                  href: `/experience/${emp.slug?.current}`,
                }))}
              />
            </Show>
          </Suspense>
        </div>
      </section>

      {/* Recent Posts Section */}
      <section class="bg-gray-800 py-16">
        <div class="container mx-auto px-4">
          <div class="mb-12 flex justify-between items-end">
            <div>
              <h2 class="text-3xl font-bold text-white mb-4">Recent Posts</h2>
              <p class="text-gray-300">Latest thoughts, tutorials, and insights</p>
            </div>
            <Link
              to="/posts"
              class="text-emerald-400 hover:text-emerald-300 transition-colors flex items-center gap-2"
            >
              View All Posts
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
              </svg>
            </Link>
          </div>

          <Suspense fallback={<div class="text-white text-center">Loading posts...</div>}>
            <Show
              when={data().posts && data().posts.length > 0}
              fallback={<div class="text-gray-400 text-center">No posts available</div>}
            >
              <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
                <For each={data().posts}>
                  {(post) => (
                    <PostCard
                      title={post.title || 'Untitled Post'}
                      description={post.excerpt}
                      bannerPath={post.banner?.asset?._ref ? `https://cdn.sanity.io/images/${post.banner.asset._ref}` : undefined}
                      createdAt={post.publishedAt || post.createdAt}
                      tags={post.tags}
                      href={`/posts/${post.slug?.current}`}
                    />
                  )}
                </For>
              </div>
            </Show>
          </Suspense>
        </div>
      </section>

      {/* Skills/Tags Section */}
      <section class="bg-gray-900 py-16">
        <div class="container mx-auto px-4">
          <div class="mb-12 flex justify-between items-end">
            <div>
              <h2 class="text-3xl font-bold text-white mb-4">Skills & Technologies</h2>
              <p class="text-gray-300">Tools and technologies I work with</p>
            </div>
            <Link
              to="/tags"
              class="text-emerald-400 hover:text-emerald-300 transition-colors flex items-center gap-2"
            >
              See All Skills/Technologies
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
              </svg>
            </Link>
          </div>

          <Suspense fallback={<div class="text-white text-center">Loading skills...</div>}>
            <Show
              when={data().tags && data().tags.length > 0}
              fallback={<div class="text-gray-400 text-center">No tags available</div>}
            >
              <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
                <For each={data().tags}>
                  {(tag) => (
                    <TagCard
                      title={tag.title || 'Untitled Tag'}
                      description={tag.description}
                      icon={tag.icon}
                      href={`/tags/${tag.slug?.current}`}
                      postCount={tag.posts?.length}
                      projectCount={tag.projects?.length}
                    />
                  )}
                </For>
              </div>
            </Show>
          </Suspense>
        </div>
      </section>

      {/* Contact/Get In Touch Section */}
      <section id="contact" class="bg-gray-800 py-16">
        <div class="container mx-auto px-4">
          <div class="mb-12 text-center">
            <h2 class="text-3xl font-bold text-white mb-4">Get In Touch</h2>
            <p class="text-gray-300">
              Interested in working together? Feel free to reach out through any of the channels below.
            </p>
          </div>

          {/* Social Links */}
          <div class="flex justify-center gap-6 mb-12">
            {/* LinkedIn */}
            <a
              href="https://www.linkedin.com/in/conner-ohnesorge"
              target="_blank"
              rel="noopener noreferrer"
              class="w-12 h-12 bg-gray-700 rounded-full flex items-center justify-center text-white hover:bg-blue-600 transition-colors"
              aria-label="LinkedIn"
            >
              <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
                <path d="M19 0h-14c-2.761 0-5 2.239-5 5v14c0 2.761 2.239 5 5 5h14c2.762 0 5-2.239 5-5v-14c0-2.761-2.238-5-5-5zm-11 19h-3v-11h3v11zm-1.5-12.268c-.966 0-1.75-.79-1.75-1.764s.784-1.764 1.75-1.764 1.75.79 1.75 1.764-.783 1.764-1.75 1.764zm13.5 12.268h-3v-5.604c0-3.368-4-3.113-4 0v5.604h-3v-11h3v1.765c1.396-2.586 7-2.777 7 2.476v6.759z" />
              </svg>
            </a>

            {/* GitHub */}
            <a
              href="https://github.com/conneroh"
              target="_blank"
              rel="noopener noreferrer"
              class="w-12 h-12 bg-gray-700 rounded-full flex items-center justify-center text-white hover:bg-gray-900 transition-colors"
              aria-label="GitHub"
            >
              <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" />
              </svg>
            </a>

            {/* Twitter/X */}
            <a
              href="https://twitter.com/conneroh"
              target="_blank"
              rel="noopener noreferrer"
              class="w-12 h-12 bg-gray-700 rounded-full flex items-center justify-center text-white hover:bg-sky-500 transition-colors"
              aria-label="Twitter"
            >
              <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
                <path d="M23 3a10.9 10.9 0 01-3.14 1.53 4.48 4.48 0 00-7.86 3v1A10.66 10.66 0 013 4s-4 9 5 13a11.64 11.64 0 01-7 2c9 5 20 0 20-11.5a4.5 4.5 0 00-.08-.83A7.72 7.72 0 0023 3z" />
              </svg>
            </a>

            {/* Email */}
            <a
              href="mailto:conner@example.com"
              class="w-12 h-12 bg-gray-700 rounded-full flex items-center justify-center text-white hover:bg-emerald-600 transition-colors"
              aria-label="Email"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
              </svg>
            </a>
          </div>

          {/* Contact Form */}
          <div class="max-w-3xl mx-auto">
            <ContactForm />
          </div>
        </div>
      </section>
    </div>
  )
}
