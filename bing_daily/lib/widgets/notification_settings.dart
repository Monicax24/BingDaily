import 'package:bing_daily/constants.dart';
import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/material.dart';
import 'package:get_storage/get_storage.dart';

/// Widget for managing notification preferences for each community.
class NotificationSettings extends StatefulWidget {
  final List<dynamic> communities;
  const NotificationSettings({super.key, required this.communities});

  @override
  State<NotificationSettings> createState() => _NotificationSettingsState();
}

class _NotificationSettingsState extends State<NotificationSettings> {
  final GetStorage _storage = GetStorage();
  final Map<String, bool> _isExpanded = {};
  final Map<String, bool> _postReminders = {};
  final Map<String, bool> _endPromptReminders = {};

  @override
  void initState() {
    super.initState();
    // Initialize storage and state for each community
    for (var community in widget.communities) {
      _isExpanded[community] = false;
      _postReminders[community] =
          _storage.read('postReminder_$community') ?? true;
      _endPromptReminders[community] =
          _storage.read('endPromptReminder_$community') ?? true;
    }
  }

  /// Toggles the expanded state of a community item.
  void _toggleExpanded(String community) {
    setState(() {
      _isExpanded[community] = !(_isExpanded[community] ?? false);
    });
  }

  /// Updates the post reminder preference and saves to storage.
  void _togglePostReminder(String community, bool value) {
    setState(() {
      _postReminders[community] = value;
      _storage.write('postReminder_$community', value);
      if (!value) {
        FirebaseMessaging.instance.unsubscribeFromTopic(
          "${community}_postReminder",
        );
      } else {
        FirebaseMessaging.instance.subscribeToTopic(
          "${community}_postReminder",
        );
      }
    });
  }

  /// Updates the end-of-prompt reminder preference and saves to storage.
  void _toggleEndPromptReminder(String community, bool value) {
    setState(() {
      _endPromptReminders[community] = value;
      _storage.write('endPromptReminder_$community', value);
    });
    if (!value) {
      FirebaseMessaging.instance.unsubscribeFromTopic(
        "${community}_endPromptReminder",
      );
    } else {
      FirebaseMessaging.instance.subscribeToTopic(
        "${community}_endPromptReminder",
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Notification Settings',
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: bingGreen,
          ),
        ),
        const SizedBox(height: 8),
        if (widget.communities.isEmpty)
          const Text(
            'No communities joined',
            style: TextStyle(fontSize: 14, color: Colors.grey),
          ),
        ...widget.communities.map((community) {
          final isExpanded = _isExpanded[community] ?? false;
          return AnimatedContainer(
            duration: const Duration(milliseconds: 200),
            height: isExpanded ? 175 : 60,
            margin: const EdgeInsets.symmetric(vertical: 4),
            decoration: BoxDecoration(
              color: Colors.grey[200],
              borderRadius: BorderRadius.circular(8),
              border: Border.all(color: Colors.grey[300]!),
            ),
            child: Column(
              children: [
                ListTile(
                  title: Text(community, style: const TextStyle(fontSize: 16)),
                  trailing: Icon(
                    isExpanded
                        ? Icons.keyboard_arrow_up
                        : Icons.keyboard_arrow_down,
                  ),
                  onTap: () => _toggleExpanded(community),
                ),
                if (isExpanded) ...[
                  SwitchListTile(
                    title: const Text(
                      'Reminder to Post',
                      style: TextStyle(fontSize: 14),
                    ),
                    value: _postReminders[community] ?? true,
                    onChanged: (value) => _togglePostReminder(community, value),
                  ),
                  SwitchListTile(
                    title: const Text(
                      'End of Prompt Reminder',
                      style: TextStyle(fontSize: 14),
                    ),
                    value: _endPromptReminders[community] ?? true,
                    onChanged: (value) =>
                        _toggleEndPromptReminder(community, value),
                  ),
                ],
              ],
            ),
          );
        }),
      ],
    );
  }
}
