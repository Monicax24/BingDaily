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
    posts: [1, 2, 3, 4, 5],
    postTime: '09:00',
    defaultPrompt: 'Share your daily campus moment',
  ),
];

// List of mock daily posts.
final List<DailyPost> mockPosts = [
  DailyPost(
    postId: 1,
    communityId: 1,
    picture: 'https://images.unsplash.com/photo-1461800919507-79b16743b257?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Mnx8aHVtYW58ZW58MHx8MHx8fDI%3D&auto=format&fit=crop&q=60&w=500',
    caption: 'Beautiful sunrise on campus',
    author: 1,
    timePosted: DateTime.now().subtract(const Duration(hours: 2)),
    likes: 5,
  ),
  DailyPost(
    postId: 2,
    communityId: 1,
    picture: 'https://images.unsplash.com/photo-1581456495146-65a71b2c8e52?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTF8fGh1bWFufGVufDB8fDB8fHwy&auto=format&fit=crop&q=60&w=500',
    caption: 'Coffee break between classes',
    author: 2,
    timePosted: DateTime.now().subtract(const Duration(hours: 1)),
    likes: 3,
  ),
  DailyPost(
    postId: 3,
    communityId: 1,
    picture: 'https://unsplash.com/photos/person-standing-on-body-of-water-4csA42uPfEo',
    caption: 'Notes from CS lecture',
    author: 3,
    timePosted: DateTime.now().subtract(const Duration(hours: 3)),
    likes: 7,
  ),
  DailyPost(
    postId: 4,
    communityId: 1,
    picture: 'https://images.unsplash.com/photo-1578916045370-25461e0cf390?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Mjh8fGh1bWFufGVufDB8fDB8fHwy&auto=format&fit=crop&q=60&w=500',
    caption: 'Morning run',
    author: 1,
    timePosted: DateTime.now().subtract(const Duration(hours: 4)),
    likes: 4,
  ),
  DailyPost(
    postId: 5,
    communityId: 1,
    picture: 'https://images.unsplash.com/photo-1503023345310-bd7c1de61c7d?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MXx8aHVtYW58ZW58MHx8MHx8fDI%3D&auto=format&fit=crop&q=60&w=500',
    caption: 'Gym session',
    author: 3,
    timePosted: DateTime.now().subtract(const Duration(minutes: 30)),
    likes: 2,
  ),
];
