// Mock data file for synthetic users, communities, and posts for testing.

import 'package:bing_daily/models/app_user.dart';
import 'package:bing_daily/models/community.dart';
import 'package:bing_daily/models/daily_post.dart';

// List of mock communities.
final List<Community> mockCommunities = [
  Community(
    communityId: 1,
    picture: 'assets/images/community1.png',
    description: 'Campus life and events',
    members: [1, 2, 3],
    moderators: [1],
    posts: [1, 2],
    postTime: '09:00',
    defaultPrompt: 'Share your daily campus moment',
  ),
  Community(
    communityId: 2,
    picture: 'assets/images/community2.png',
    description: 'Study tips and academics',
    members: [2, 3],
    moderators: [2],
    posts: [3],
    postTime: '10:00',
    defaultPrompt: 'What are you studying today?',
  ),
  Community(
    communityId: 3,
    picture: 'assets/images/community3.png',
    description: 'Sports and fitness',
    members: [1, 3],
    moderators: [3],
    posts: [4, 5],
    postTime: '08:00',
    defaultPrompt: 'Show your workout routine',
  ),
];

// List of mock daily posts.
final List<DailyPost> mockPosts = [
  DailyPost(
    postId: 1,
    communityId: 1,
    picture: 'assets/images/post1.jpg',
    caption: 'Beautiful sunrise on campus',
    author: 1,
    timePosted: DateTime.now().subtract(const Duration(hours: 2)),
    likes: 5,
  ),
  DailyPost(
    postId: 2,
    communityId: 1,
    picture: 'assets/images/post2.jpg',
    caption: 'Coffee break between classes',
    author: 2,
    timePosted: DateTime.now().subtract(const Duration(hours: 1)),
    likes: 3,
  ),
  DailyPost(
    postId: 3,
    communityId: 2,
    picture: 'assets/images/post3.jpg',
    caption: 'Notes from CS lecture',
    author: 3,
    timePosted: DateTime.now().subtract(const Duration(hours: 3)),
    likes: 7,
  ),
  DailyPost(
    postId: 4,
    communityId: 3,
    picture: 'assets/images/post4.jpg',
    caption: 'Morning run',
    author: 1,
    timePosted: DateTime.now().subtract(const Duration(hours: 4)),
    likes: 4,
  ),
  DailyPost(
    postId: 5,
    communityId: 3,
    picture: 'assets/images/post5.jpg',
    caption: 'Gym session',
    author: 3,
    timePosted: DateTime.now().subtract(const Duration(minutes: 30)),
    likes: 2,
  ),
];

// List of mock users.
final List<AppUser> mockUsers = [
  AppUser(
    userId: '1',
    email: 'user1@binghamton.edu',
    username: 'userone',
    joinDate: '2025-01-01',
    communities: [1, 3],
  ),
  AppUser(
    userId: '2',
    email: 'user2@binghamton.edu',
    username: 'usertwo',
    joinDate: '2025-02-01',
    communities: [1, 2],
  ),
  AppUser(
    userId: '3',
    email: 'user3@binghamton.edu',
    username: 'userthree',
    joinDate: '2025-03-01',
    communities: [2, 3],
  ),
];
