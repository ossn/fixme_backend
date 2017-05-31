from rest_framework import serializers
from core.models import UserRepo


class UserRepoSerializer(serializers.ModelSerializer):
    class Meta:
        model = UserRepo
        fields = ('id', 'user', 'repo')