import 'package:bing_daily/constants.dart';
import 'package:bing_daily/models/daily_post.dart';
import 'package:flutter/material.dart';

/// Widget to display a single daily post with image, caption, and likes.
class PostCard extends StatelessWidget {
  final DailyPost post;

  const PostCard({super.key, required this.post});

  @override
  Widget build(BuildContext context) {
    return Card(
      color: bingWhite,
      margin: const EdgeInsets.symmetric(vertical: 8, horizontal: 16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Image.network(
            post.picture,
            height: 200,
            width: double.infinity,
            fit: BoxFit.cover,
            errorBuilder: (context, error, stackTrace) =>
                const Icon(Icons.error),
          ),
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: Text(
              post.caption,
              style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
            ),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8.0),
            child: Row(
              children: [
                const Icon(Icons.thumb_up, size: 16, color: bingGreen),
                const SizedBox(width: 4),
                Text('${post.likes} likes'),
              ],
            ),
          ),
          const SizedBox(height: 8),
        ],
      ),
    );
  }
}
