from core.models import Repository, Issue, Project
from core.serializers import RepositorySerializer, IssueSerializer, ProjectSerializer
from rest_framework import generics
from rest_framework.filters import OrderingFilter
from rest_framework.views import APIView
from rest_framework.renderers import JSONRenderer, BrowsableAPIRenderer
from rest_framework.response import Response
from rest_framework import viewsets
from django_filters import Filter, FilterSet
from django_filters.rest_framework import DjangoFilterBackend
from django.db.models import Count


def generateFilterValues(value):
    values = []
    str = ""
    for i in value:
        if i in ["[", "\"", "\\"]:
            continue
        if i in [",", "]"]:
            values.append(str)
            str = ""
            continue
        str += i
    return values


class ListFilter(Filter):
    def filter(self, qs, value):
        if not value:
            return qs

        self.lookup_expr = 'in'
        return super(ListFilter, self).filter(qs, generateFilterValues(value))


class IssueFilter(FilterSet):
    language = ListFilter(name='language')
    tech_stack = ListFilter(name='tech_stack')
    experience_needed = ListFilter(name='experience_needed')
    expected_time = ListFilter(name='expected_time')
    issue_type = ListFilter(name='issue_type')

    class Meta:
        model = Issue
        fields = ('language', 'tech_stack',
                  'experience_needed', 'expected_time', 'issue_type')


class RepositoryList(generics.ListAPIView):
    """
    Returns a list of issues, by optionally filtering against
    `repo` and `user` query parameter in the URL.
    """
    queryset = Repository.objects.all()
    serializer_class = RepositorySerializer
    filter_backends = (DjangoFilterBackend,)
    filter_fields = ('repository_url',)


class IssueList(generics.ListAPIView):
    """
    Returns a list of issues, by optionally filtering against
    `language`, `tech_stack`, `experience_needed` and `expected_time`
    query parameter in the URL.
    """
    queryset = Issue.objects.all()
    serializer_class = IssueSerializer
    filter_backends = (DjangoFilterBackend, OrderingFilter)
    filter_class = IssueFilter
    ordering_fields = ('experience_needed', 'expected_time')
    filter_backends = (DjangoFilterBackend, OrderingFilter)


class ProjectList(generics.ListAPIView):
    """
    Returns a list of projects
    """
    queryset = Project.objects.all()
    serializer_class = ProjectSerializer
    filter_backends = (DjangoFilterBackend, OrderingFilter)
    filter_fields = ("id", "display_name")
    ordering_fields = ('display_name')
    filter_backends = (DjangoFilterBackend, OrderingFilter)


class MetaData(APIView):
    """
    Returns a list of all the `language`, `tech_stack`,
    `experience_needed` present in the database.
    """
    renderer_classes = (JSONRenderer, BrowsableAPIRenderer,)

    def get(self, request, format=None):
        language = 'language'
        tech_stack = 'tech_stack'
        experience_needed = 'experience_needed'

        queryset = Issue.objects.values(
            language, tech_stack, experience_needed,)
        meta_data = {}
        meta_data[language] = set([tup[language] for tup in queryset])
        meta_data[tech_stack] = set([tup[tech_stack] for tup in queryset])
        meta_data[experience_needed] = set(
            [tup[experience_needed] for tup in queryset])
        return Response(meta_data)
